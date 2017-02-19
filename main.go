package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "./example/example.go", nil, 0) // parser.ParseComments
	if err != nil {
		log.Panic(err)
	}

	// ast.Fprint(os.Stdout, fset, node, nil)

	ast.Inspect(node, func(n ast.Node) bool {
		fcExp, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if fcExp.Name.Name == "main" {
			return true
		}

		r := fcExp.Type.Results
		if r == nil {
			return true
		}

		var hasReturnError bool
		for _, e := range r.List {
			tp, ok := e.Type.(*ast.Ident)
			if !ok {
				continue
			}

			if tp.Name == "error" { // || tp.Name == "bool" {
				hasReturnError = true
				break
			}
		}
		if !hasReturnError {
			return true
		}

		body := fcExp.Body
		if body == nil {
			return true
		}

		log.Printf("%#v", body.List)

		return true
	})
}
