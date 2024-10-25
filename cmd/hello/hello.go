package main

import "github.com/miekg/bpf"

func main() {
	bpf.TracePrintk("Hello world!\n")
}
