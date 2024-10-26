package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"slices"
	"text/template"

	"github.com/miekg/bpf"
)

var fOut = flag.String("out", "", "file name to generate")

func main() {
	flag.Parse()

	fset := token.NewFileSet()
	prog, err := parser.ParseFile(fset, flag.Arg(0), nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	ctx := bpf.New()
	ast.Walk(ctx, prog)

	slices.Reverse(ctx.Insns) // somewhat naive... ?
	tmpl, err := template.New("out").Parse(out)
	if err != nil {
		log.Fatal(err)
	}
	buf := &bytes.Buffer{}

	if err := tmpl.Execute(buf, ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)
}
