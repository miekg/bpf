package main

import "go/ast"

type BasicLit struct {
	Key int // "numeric" index into the map where this basiclit is stored.
	ast.BasicLit
}
