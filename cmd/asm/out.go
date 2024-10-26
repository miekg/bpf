package main

const out = `// Code generated by "asm"; DO NOT EDIT.
package main

import (
        "log"

        "github.com/miekg/bpf"
        "github.com/cilium/ebpf"
        "github.com/cilium/ebpf/asm"
        "github.com/cilium/ebpf/link"
//        "github.com/cilium/ebpf/perf"
        "github.com/cilium/ebpf/rlimit"
)

var ctx = bpf.New()

var progSpec = &ebpf.ProgramSpec{
	Name:    "my_trace_prog",
	Type:    ebpf.TracePoint,
	License: "GPL",
}

func main() {
	setupRLimit()

	setupROMap()

	setupProg()

        prog, err := ebpf.NewProgram(progSpec)
        if err != nil {
		log.Fatalf("creating ebpf program: %s", err)
        }
        defer prog.Close()

        tp, err := link.Tracepoint("syscalls", "sys_enter_openat", prog, nil)
        if err != nil {
		log.Fatalf("opening tracepoint: %s", err)
        }
        defer tp.Close()
}

func setupROMap() {
	{{- range $i,$v := .Constants }}
	ctx.ROdata.Put({{$i}}, {{$v.Value}})
	{{end}}
}

func setupProg() {
	progSpec.Instructions = asm.Instructions{
	{{- range .Insns}}
		{{. -}},
	{{- end}}
	}
}

func setupRLimit() {
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}
}
`
