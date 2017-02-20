package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func main() {
	file := flag.String("file", "", "Filename to analyze")
	dir := flag.String("dir", "", "Directory to analyze")

	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}
	var err error
	var node ast.Node

	fset := token.NewFileSet()
	switch {
	case *file != "":
		node, err = parser.ParseFile(fset, *file, nil, 0) // parser.ParseComments
		if err != nil {
			panic(err)
		}
	case *dir != "":
		// pkgs, err = parser.ParseDir(fset, *dir, nil, 0)
		fmt.Println("TODO")
		return
	default:
		return
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
		// fmt.Println(fcExp.Name.Name)

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
			in, rn, hasOHR := analizeErrReturn(fcExp.Body)
			if hasOHR {
				ifPos := fset.Position(in.If)
				fmt.Printf("\n%s:%d\n\n", ifPos.Filename, ifPos.Line)

				printer.Fprint(os.Stdout, fset, in)
				fmt.Println()
				printer.Fprint(os.Stdout, fset, rn)
				fmt.Printf("\n---------------------------\n")
			}
		}

		return true
	})
}

//
func analizeErrReturn(body *ast.BlockStmt) (*ast.IfStmt, *ast.ReturnStmt, bool) {
	if body == nil {
		return nil, nil, false
	}

	lenBL := len(body.List)
	// we're trying to detect at least 2 lines: 1 if and 1 return
	if lenBL < 2 {
		return nil, nil, false
	}

	ifSt, ok := body.List[lenBL-2].(*ast.IfStmt)
	if !ok {
		return nil, nil, false
	}

	if !processIfNode(ifSt) {
		return nil, nil, false
	}

	retSt, ok := body.List[lenBL-1].(*ast.ReturnStmt)
	if !ok {
		return nil, nil, false
	}

	if len(retSt.Results) != 1 {
		return nil, nil, false
	}
	// fmt.Printf("%#v", retSt.Results[0]) // return statement

	return ifSt, retSt, true
}

func processIfNode(ifSt *ast.IfStmt) bool {
	binCond, ok := ifSt.Cond.(*ast.BinaryExpr)
	if !ok {
		// fmt.Printf("%#v", ifSt.Cond)
		// check for bool with UniaryExpr
		_, ok = ifSt.Cond.(*ast.Ident)
		if !ok {
			return false
		}
		return true
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

	// check return inside of 'if'
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

	if varName == ident.Name || ident.Name == "true" || ident.Name == "false" {
		// fmt.Printf("%#v", retSt.Results[0]) // return statement
		return true
	}

	return false
}
