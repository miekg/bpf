package bpf

import "github.com/cilium/ebpf"

// Context is an BPF context, it holds (global) data that each bpf Go program needs to access while running.
type Context struct {
	ROdata *ebpf.Map // RO is a map that holds the .rodata for a eBPF program.
}

// New returns a new context.
func New() *Context {
	c := &Context{}

	rodata, _ := ebpf.NewMap(&ebpf.MapSpec{
		Type: ebpf.Array,
		Name: "rodata",
	})
	c.ROdata = rodata

	return c
}
