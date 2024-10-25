package ebpf

// TracePrintk is bpf_trace_printk, see https://docs.ebpf.io/linux/helper-function/bpf_trace_printk/ .
func TracePrintk(format string, a ...any) int { return 0 }
