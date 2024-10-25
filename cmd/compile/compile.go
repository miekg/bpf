package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	flag.Parse()

	fset := token.NewFileSet()
	prog, err := parser.ParseFile(fset, flag.Arg(0), nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(prog, func(n ast.Node) bool {
		funcCall, ok := n.(*ast.CallExpr)
		if ok {
			fmt.Println(funcCall.Fun)
		}
		return true
	})
	ast.Print(fset, prog)
}
