package main

import (
	"github.com/miekg/ebpf"
)

func main() {
	ebpf.TracePrintk("Hello world!\n")
}
