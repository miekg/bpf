package main

import (
	"embed"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/miekg/bpf/context"
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

	ast.Inspect(prog, func(n ast.Node) bool {
		funcCall, ok := n.(*ast.CallExpr)
		if ok {
			fmt.Println(funcCall.Fun)
		}
		return true
	})
	ast.Print(fset, prog)

	ctx := context.New()
	ast.Walk(ctx, prog)
}
