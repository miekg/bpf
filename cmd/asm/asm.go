package main

import (
	"embed"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"slices"

	"github.com/miekg/bpf"
)

//go:embed *.tmpl
var tmplfs embed.FS

/*
 if tmpl, err = tmpl.ParseFS(fs, path.Join(cmdline.Pkg(reflect.TypeOf(m)), file)); err != nil {
                        log.Fatalf("Failed to generate manual page: %s", err)
                }
*/

var fOut = flag.String("out", "", "file name to generate")

func main() {
	flag.Parse()

	fset := token.NewFileSet()
	prog, err := parser.ParseFile(fset, flag.Arg(0), nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	ast.Print(fset, prog)

	ctx := bpf.New()
	ast.Walk(ctx, prog)

	println(len(ctx.Insns))
	slices.Reverse(ctx.Insns) // somewhat naive... ?
	for _, s := range ctx.Insns {
		fmt.Println(s)
	}
}
