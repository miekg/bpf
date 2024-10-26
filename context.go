package bpf

import (
	"go/ast"
	"log"
	"strconv"

	"github.com/cilium/ebpf"
)

// Context is an BPF context, it holds (global) data that each bpf Go program needs to access while running. It is used
// when parsing a Go program and then generating the new ebpf program.
type Context struct {
	ROdata *ebpf.Map // RO is a map that holds the .rodata for a eBPF program, used when generating.

	Constants []*ast.BasicLit // Constants holds all the constants.
	Insns     []string        // Insns are the assembly instruction in Go syntax.
}

// New returns a new context.
func New() *Context {
	c := &Context{}
	c.Insns = []string{
		// set the error code for the ebpf to zero, in reverse order because we slices.Reverse this (TODO(xxx))
		`asm.Return()`,
		`asm.Mov.Imm(asm.R0, 0)`,
	}
	return c
}

func (ctx *Context) Map() {
	rodata, err := ebpf.NewMap(&ebpf.MapSpec{
		Type:       ebpf.Array,
		KeySize:    4,
		ValueSize:  8,
		Name:       "rodata",
		MaxEntries: 10,
	})
	if err != nil {
		log.Fatal(err)
	}
	ctx.ROdata = rodata
}

// FD returns the fd of the rodata map.
func (ctx *Context) FD() int { return ctx.ROdata.FD() }

func (ctx *Context) AddIns(ins string) { ctx.Insns = append(ctx.Insns, ins) }

// AddConstant adds a constant to the context, it returns the index of the element so the caller can use that as a reference.
func (ctx *Context) AddConstant(basicLit *ast.BasicLit) string {
	ctx.Constants = append(ctx.Constants, basicLit)
	return strconv.FormatUint(uint64(len(ctx.Constants)-1), 10) // zero based, hence -1
}

// Visitor function to traverse AST and generate code
func (ctx *Context) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.ImportSpec: // skip any imports
		return nil
	case *ast.IfStmt:
		ctx.genIfStmt(n)
	case *ast.AssignStmt:
		ctx.genAssignStmt(n)
	case *ast.FuncDecl:
		ctx.genFuncDecl(n)
	case *ast.CallExpr:
		ctx.genCallExpr(n)
	case *ast.BasicLit:
		ctx.genBasicLit(n)
	}
	return ctx
}

func (ctx *Context) genCallExpr(callExpr *ast.CallExpr) {
	if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		switch fun.Sel.Name {
		case "TracePrintk":
			ctx.AddIns(`asm.FnTracePrintk.Call()`)
		}
	}
}

func (ctx *Context) genFuncDecl(funcDecl *ast.FuncDecl) {
	// check for builtins, or there others??
	// check package??
	// println(funcDecl.Name.Name)

	/*
		asm.LoadMapPtr(asm.R2, events.FD()), // file descriptor of the perf event array
			asm.LoadImm(asm.R3, 0xffffffff, asm.DWord),
			asm.Mov.Reg(asm.R4, asm.RFP),
			asm.Add.Imm(asm.R4, -8),
			asm.Mov.Imm(asm.R5, 4),

			// call FnPerfEventOutput, an eBPF kernel helper
			asm.FnPerfEventOutput.Call()
	*/
}

func (ctx *Context) genIfStmt(ifStmt *ast.IfStmt) {}

func (ctx *Context) genAssignStmt(assignStmt *ast.AssignStmt) {}

func (ctx *Context) genBasicLit(basicLit *ast.BasicLit) {
	index := ctx.AddConstant(basicLit) // order in generated program will be the same
	ctx.AddIns(`asm.LoadMapValue(asm.R1, ctx.FD(), ` + index + `)`)
}
