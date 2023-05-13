package parser

import (
	"go/ast"
	"strings"

	"honnef.co/go/tools/go/ast/astutil"
)

type replceStmt interface {
	replace(c *astutil.Cursor)
}

func parseAstFile(p *SourcePackage, file *ast.File) {
	stmts := make(map[ast.Node]replceStmt)

	ast.Inspect(file, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncDecl:
			findReplacementExpr(p, stmts, n)
			return false
		}

		return true
	})

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		cn := c.Node()

		if stmt, ok := stmts[cn]; ok {
			stmt.replace(c)
		}

		return true
	}, nil)

	// replace imports path
	for _, x := range file.Imports {
		if strings.HasPrefix(x.Path.Value, "\""+p.pkg.Module.Path) {
			x.Path.Value = strings.Replace(x.Path.Value, p.pkg.Module.Path, p.pkg.Module.Path+"/dist", 1)
		}
	}
}

func findReplacementExpr(p *SourcePackage, stmts map[ast.Node]replceStmt, n ast.Node) {
	var parentNode ast.Node

	fnScope := newFunctionScope(p, n)
	ast.Inspect(n, func(cn ast.Node) bool {
		if n == cn {
			return true
		}

		switch x := cn.(type) {
		case *ast.FuncLit:
			findReplacementExpr(p, stmts, cn)
			return false

		case *ast.ReturnStmt:
			return false

		case *ast.CallExpr:
			if fresults, ok := hasFuncResultsError(p, x); ok {
				switch pn := (parentNode).(type) {
				case *ast.AssignStmt:
					if pn.Rhs[0] == x {
						if lhs, needReplace := normolizeAssignStmt(p, pn.Lhs, fresults, fnScope.getNextErrorName()); needReplace {
							stmts[parentNode] = &replceCallExprStmt{
								parentNode: parentNode,
								callExpr:   x,
								lhs:        lhs,
								fnScope:    fnScope,
							}
						}
					}

				case *ast.ExprStmt:
					if pn.X == x {
						if lhs, needReplace := normolizeAssignStmt(p, []ast.Expr{}, fresults, fnScope.getNextErrorName()); needReplace {
							stmts[parentNode] = &replceCallExprStmt{
								parentNode: parentNode,
								callExpr:   x,
								lhs:        lhs,
								fnScope:    fnScope,
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
