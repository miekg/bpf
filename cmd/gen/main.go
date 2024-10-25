package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	fset := token.NewFileSet()
	prog, err := parser.ParseFile(fset, "hello.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Inspect(prog, func(n ast.Node) bool {
		// Find Function Call Statements
		funcCall, ok := n.(*ast.CallExpr)
		if ok {
			fmt.Println(funcCall.Fun)
		}
		return true
	})
	ast.Print(fset, prog)

	// Create a loader.Config to hold the program configuration
	var conf loader.Config
	conf.Fset = fset
	conf.CreateFromFiles("main", prog)

	// Load the program from the parsed file
	prog2, err := conf.Load()
	if err != nil {
		fmt.Println("Error loading program:", err)
		return
	}

	// Create an SSA program from the loaded program
	ssaprog := ssautil.CreateProgram(prog2, ssa.SanityCheckFunctions)

	// Build the SSA representation for all packages in the program
	ssaprog.Build()

	// Print the SSA for each function in the "main" package
	for _, pkg := range ssaprog.AllPackages() {
		fmt.Printf("%+v\n", pkg)
		b := &bytes.Buffer{}
		ssa.WritePackage(b, pkg)
		fmt.Printf("%s\n", b)
		for _, member := range pkg.Members {
			if fn, ok := member.(*ssa.Function); ok {
				b.Reset()
				fmt.Printf("Function %s SSA:\n", fn.Name())
				fmt.Printf("+%v\n", fn)
				ssa.WriteFunction(b, fn)
				fmt.Printf("%s\n", b)
			}
		}
	}

}
