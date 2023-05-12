package parser2

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"honnef.co/go/tools/go/ast/astutil"
)

func parseAstFile(p *SourcePackage, file *ast.File) {
	astutil.Apply(file, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch n.(type) {
		case *ast.FuncDecl:
			c.Replace(replaceErrors(p, &n))
		}

		return true
	}, nil)

	// astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
	// 	n := c.Node()
	// 	switch x := n.(type) {
	// 	case *ast.CallExpr:
	// 		fmt.Println(x)
	// 		switch s := x.Fun.(type) {
	// 		case *ast.Ident:
	// 			fmt.Print(s.Name, s.Obj)
	// 		case *ast.SelectorExpr:
	// 			obj := p.pkg.TypesInfo.Uses[s.Sel]
	// 			fmt.Println("sel", s.Sel, p.pkg.TypesInfo.Uses[s.Sel])
	// 			if f, ok := obj.(*types.Func); ok {
	// 				ss := f.Type().(*types.Signature)
	// 				fmt.Println(f.Type(), ss.Results().At(1).String(), ss.Results().At(1).Type().String() == "error")
	// 				// fmt.Println(f.Type(), reflect.TypeOf(f.Type()).NumOut())
	// 			}
	// 		}
	// 	}

	// 	return true
	// })
}

func replaceErrors(p *SourcePackage, n *ast.Node) ast.Node {
	var parentNode *ast.Node

	funcWrapper := newFuncWrapper(p, n)
	errInc := 0
	return astutil.Apply(*n, func(c *astutil.Cursor) bool {
		cn := c.Node()
		parentNode = &cn

		if *n == cn {
			return true
		}

		switch x := cn.(type) {
		case *ast.SelectStmt:
			id, ok := x.Fun.(*ast.Ident)

			if ok {
				if id.Name == "funcWithOneError" {
					c.Replace(&ast.UnaryExpr{
						Op: token.NOT,
						X:  x,
					})
				}
			}
		case *ast.FuncLit:
			c.Replace(replaceErrors(p, &cn))
			return false
		case *ast.ExprStmt:
			if c := x.X.(*ast.CallExpr); c != nil {
				if hasFuncError(p, c) {
					fmt.Println("***************fds fsd fsdf sdf sdf sdf sdfsd ")
				}
			}
		case *ast.AssignStmt:
			l := x.Lhs[len(x.Lhs)-1]
			v, ok := l.(*ast.Ident)

			if ok && v.Name == "_" {

				errInc++
				v.Name = "err_" + strconv.Itoa(errInc)

				c.InsertAfter(&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						// err
						X: &ast.Ident{Name: v.Name},
						// !=
						Op: token.NEQ,
						// nil
						Y: &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							getReturnErrorStm(funcWrapper, v.Name),
						},
					},
				})
			}

		// exit if nested func declaration
		case *ast.FuncDecl:
			return false
		}

		return true
	}, nil)
}

func getReturnErrorStm(fw *funcWrapper, nameErr string) ast.Stmt {
	fmt.Println(fw.getName())

	if !fw.hasErrorResults() {
		panic(fmt.Errorf("func %s not return error", fw.getName()))
	}

	results := make([]ast.Expr, 0)

	for _, v := range fw.getResults().List {
		switch x := v.Type.(type) {
		case *ast.StarExpr:
			results = append(results, &ast.Ident{Name: "nil"})
		case *ast.Ident:
			if x.Name == "error" {
				results = append(results, &ast.Ident{Name: nameErr})
			} else {
				results = append(results, &ast.Ident{Name: getDefaultValue(x.Name)})
			}
		default:
			results = append(results, &ast.Ident{Name: "nil"})
		}
	}

	return &ast.ReturnStmt{
		Results: results,
	}
}
