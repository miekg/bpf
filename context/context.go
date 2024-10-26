package context

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/asm"
)

// Context is an BPF context, it holds (global) data that each bpf Go program needs to access while running.
type Context struct {
	ROdata *ebpf.Map        // RO is a map that holds the .rodata for a eBPF program.
	Insns  asm.Instructions // Insns are the assembly instruction in Go syntax.
}

// New returns a new context.
func New() *Context {
	c := &Context{Insns: asm.Instructions{}}

	rodata, _ := ebpf.NewMap(&ebpf.MapSpec{
		Type: ebpf.Array,
		Name: "rodata",
	})
	c.ROdata = rodata

	return c
}

// Visitor function to traverse AST and generate code
func (ctx *Context) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.FuncDecl:
		ctx.genFuncDecl(n)
	case *ast.IfStmt:
		ctx.genIfStmt(n)
	case *ast.AssignStmt:
		ctx.genAssignStmt(n)
	case *ast.BasicLit:
		ctx.genBasicLit(n)
	}
	return ctx
}

func (ctx *Context) genFuncDecl(funcDecl *ast.FuncDecl) {
	ctx.code = append(ctx.code, fmt.Sprintf("func_%s:", funcDecl.Name.Name))
}

func (ctx *Context) genIfStmt(ifStmt *ast.IfStmt) {
	ctx.code = append(ctx.code, "cmp ...")
	ctx.code = append(ctx.code, "jne ...")
}

func (ctx *Context) genAssignStmt(assignStmt *ast.AssignStmt) {
	ctx.code = append(ctx.code, "mov ...")
}

func (ctx *Context) genBasicLit(basicLit *ast.BasicLit) {
	ctx.code = append(ctx.code, "mov ...")
}

// Main function to start the AST traversal
func main() {
	src := `
        package main
        func main() {
            if x > 0 {
                x = 1
            }
        }
    `
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := New()
	ast.Walk(ctx, node)

	for _, line := range ctx.code {
		fmt.Println(line)
	}
}
