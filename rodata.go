package ebpf

import (
	"github.com/cilium/ebpf"
)

func RODATA() error {
	// This rodata map holds all constants fr the program.
	rodata, err := ebpf.NewMap(&ebpf.MapSpec{
		Type: ebpf.Array,
		Name: "rodata",
	})
	if err != nil {
		return err
	}
	defer rodata.Close()
	return nil
}

r1 = map? FD() from map
r2 = element
FnMapLookupElem.Call()
