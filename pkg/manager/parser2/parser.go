package parser2

import (
	"fmt"
	"go/ast"
	"go/token"

	"honnef.co/go/tools/go/ast/astutil"
)

type replceStmt struct {
	parentNode ast.Node
	callExpr   *ast.CallExpr
	fw         *funcWrapper
	lhs        []ast.Expr
}

func parseAstFile(p *SourcePackage, file *ast.File) {
	stmts := make(map[ast.Node]*replceStmt)

	ast.Inspect(file, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncDecl:
			replaceErrors(p, stmts, n)
		}

		return true
	})

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		cn := c.Node()

		if stmt, ok := stmts[cn]; ok {
			replaceReturnErrorStmt(c, stmt.fw, stmt.lhs, stmt.callExpr)
		}

		return true
	}, nil)
}

func replaceErrors(p *SourcePackage, stmts map[ast.Node]*replceStmt, n ast.Node) {
	var parentNode ast.Node

	funcWrapper := newFuncWrapper(p, n)
	ast.Inspect(n, func(cn ast.Node) bool {
		if n == cn {
			return true
		}

		switch x := cn.(type) {
		case *ast.FuncLit:
			replaceErrors(p, stmts, cn)
			return false
		case *ast.CallExpr:
			if fresults, ok := hasFuncResultsError(p, x); ok {
				switch pn := (parentNode).(type) {
				case *ast.AssignStmt:
					if pn.Rhs[0] == x {
						if lhs, ok := normolizeAssignStmt(p, pn.Lhs, fresults, funcWrapper.getNextErrorName()); ok {
							stmts[parentNode] = &replceStmt{
								parentNode: parentNode,
								callExpr:   x,
								lhs:        lhs,
								fw:         funcWrapper,
							}
						}
					}
				case *ast.ExprStmt:
					if pn.X == x {
						if lhs, ok := normolizeAssignStmt(p, []ast.Expr{}, fresults, funcWrapper.getNextErrorName()); ok {
							stmts[parentNode] = &replceStmt{
								parentNode: parentNode,
								callExpr:   x,
								lhs:        lhs,
								fw:         funcWrapper,
							}
						}
					}
				}
			}

		// exit if nested func declaration
		case *ast.FuncDecl:
			return false

		case *ast.AssignStmt, *ast.ExprStmt:
			parentNode = cn
		}

		return true
	})
}

func replaceReturnErrorStmt(c *astutil.Cursor, fw *funcWrapper, lhs []ast.Expr, ce *ast.CallExpr) {
	if !fw.hasErrorResults() {
		panic(fmt.Errorf("func %s not return error", fw.getName()))
	}

	nameErr := lhs[len(lhs)-1].(*ast.Ident).Name
	c.Replace(&ast.AssignStmt{
		Lhs: lhs,
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			ce,
		},
	})

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

	c.InsertAfter(&ast.IfStmt{
		Cond: &ast.BinaryExpr{
			// err
			X: &ast.Ident{Name: nameErr},
			// !=
			Op: token.NEQ,
			// nil
			Y: &ast.Ident{Name: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: results,
				},
			},
		},
	})
}
