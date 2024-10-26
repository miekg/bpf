helloworld-ebpf.go: helloworld.go
	../../cmd/asm/asm helloworld.go > helloworld-ebpf.go

PHONY: build
build:
	go build -o helloworld-ebpf helloworld-ebpf.go
