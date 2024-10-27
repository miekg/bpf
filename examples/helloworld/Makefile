.PHONY: helloworld-ebpf.go
helloworld-ebpf.go:
	../../cmd/asm/asm helloworld.go > helloworld-ebpf.go

.PHONY: build
build:
	go build -o helloworld-ebpf helloworld-ebpf.go
