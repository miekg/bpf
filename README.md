# (e)BPF in *pure* Go

> eBPF in C? What am I, a farmer?

The eBPF ecosystem heavily depends on clang/llvm to compile (C) source code into object code that
can be loaded by the kernel. The goal here is to get rid of clang/llvm and use pure Go.

The benefit of this approach is that you can use _all_ the Go development tooling for writing an eBPF
program. Cilium has been doing _a lot_ of work in this space.

This package is `bpf` all Cilium stuff is `ebpf`.

When using eBPF you can also talk with it via FDs and eBPF maps, this means the non-eBPF bit of your
code also lives somewhere. Taking this all into account I've finally settled on the following
approach.

Using Cilium' asm package you can already do things like this ([example from](https://)).

~~~ go
// Minimal program that writes the static value '123' to the perf ring on
// each event. Note that this program refers to the file descriptor of
// the perf event array created above, which needs to be created prior to the
// program being verified by and inserted into the kernel.
progSpec.Instructions = asm.Instructions{
	// store the integer 123 at FP[-8]
	asm.Mov.Imm(asm.R2, 123),
	asm.StoreMem(asm.RFP, -8, asm.R2, asm.Word),

	// load registers with arguments for call of FnPerfEventOutput
	asm.LoadMapPtr(asm.R2, events.FD()), // file descriptor of the perf event array
	asm.LoadImm(asm.R3, 0xffffffff, asm.DWord),
	asm.Mov.Reg(asm.R4, asm.RFP),
	asm.Add.Imm(asm.R4, -8),
	asm.Mov.Imm(asm.R5, 4),

	// call FnPerfEventOutput, an eBPF kernel helper
	asm.FnPerfEventOutput.Call(),

	// set exit code to 0
	asm.Mov.Imm(asm.R0, 0),
	asm.Return(),
}
~~~

Which is eBPF assembly with Go functions and types - meaning Cilium already made a assembler.

What if we can generate the above from Go code using Go code. This is essentially a eBPF assembly in
Go and then piggybacking on all the Cilium stuff.

Thus the above assembly _should_ be generated from the following Go code:

~~~ go
func BPF() {
    gobpf.PerfEventOutput(bpf.Context, events.FD(), gobpf.BPF_F_INDEX_MASK, 123)
  }
~~~

Where `gobpf` is _this_ library and `bpf` is cilium's. As you can see `event.FD()` is undeclared, so
it might make sense to generate the whole program via some `bpf` comment tags or some other mechanism.

## TODO

- constants on the stack
- generate complete program, looks like a lot of boiler plate

## Stuff of interest

* go-delve/delve@v1.22.1/pkg/proc/internal/ebpf/helpers.go
* github.com:miekg/ebpfcat (forked)
* https://github.com/DQNEO/babygo

## Reading list

- https://qmonnet.github.io/whirl-offload/2020/04/12/llvm-ebpf-asm/
- https://pkg.go.dev/golang.org/x/tools/cmd/ssadump
- https://go.googlesource.com/tools/+/master/go/ssa/example_test.go
- https://benhoyt.com/writings/mugo/
