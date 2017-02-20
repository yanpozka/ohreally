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

		// ignore main function
		if fcExp.Name.Name == "main" {
			return true
		}
		// log.Println(fcExp.Name.Name)

		// check if current function returns something
		r := fcExp.Type.Results
		if r == nil {
			return true
		}

		// only if returns one element (bool or error)
		if len(r.List) != 1 {
			return true
		}

		tp, ok := r.List[0].Type.(*ast.Ident)
		if !ok {
			return true
		}

		if tp.Name == "error" || tp.Name == "bool" {
			analizeErrReturn(fcExp.Body)
		}

		return true
	})
}

// TODO: return if position and return position
func analizeErrReturn(body *ast.BlockStmt) bool {
	if body == nil {
		return false
	}

	lenBL := len(body.List)
	// we're trying to detect a least 2 lines: 1 if and 1 return
	if lenBL < 2 {
		return false
	}

	ifSt, ok := body.List[lenBL-2].(*ast.IfStmt)
	if !ok {
		return false
	}

	{
		binCond, ok := ifSt.Cond.(*ast.BinaryExpr)
		if !ok {
			// check for bool with UniaryExpr
			return false
		}

		// check == or != comparation
		if (binCond.Op == token.EQL || binCond.Op == token.NEQ) == false {
			return false
		}

		var compWithNil bool
		var varName string

		leftOp, ok := binCond.X.(*ast.Ident)
		if !ok {
			return false
		}
		if leftOp.Name == "nil" {
			compWithNil = true
		} else {
			varName = leftOp.Name
		}

		rightOp, ok := binCond.Y.(*ast.Ident)
		if !ok {
			return false
		}
		if rightOp.Name == "nil" {
			compWithNil = true
		} else {
			varName = rightOp.Name
		}

		if !compWithNil {
			return false
		}

		ifBody := ifSt.Body
		if ifBody == nil || len(ifBody.List) != 1 {
			return false
		}

		retSt, ok := ifBody.List[0].(*ast.ReturnStmt)
		if !ok {
			return false
		}
		if len(retSt.Results) != 1 {
			return false
		}
		ident, ok := retSt.Results[0].(*ast.Ident)
		if !ok {
			return false
		}

		if varName == ident.Name {
			log.Printf("%#v", retSt.Results[0]) // return statement
		} else {
			return false
		}
	}

	retSt, ok := body.List[lenBL-1].(*ast.ReturnStmt)
	if !ok {
		return false
	}

	if len(retSt.Results) != 1 {
		return false
	}
	log.Printf("%#v", retSt.Results[0]) // return statement

	return true
}
