package parser

import (
	"go/ast"
	"go/printer"
	"os"
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

		defer func() {
			if r := recover(); r != nil {
				if err := printer.Fprint(os.Stdout, p.pkg.Fset, cn); err != nil {
					panic(err)
				}

				panic(r)
			}
		}()

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
	fnScope := newFunctionScope(p, n)
	ast.Inspect(n, func(cn ast.Node) bool {
		defer func() {
			if r := recover(); r != nil {
				if err := printer.Fprint(os.Stdout, p.pkg.Fset, cn); err != nil {
					panic(err)
				}
				panic(r)
			}
		}()

		if n == cn {
			return true
		}

		switch x := cn.(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			findReplacementExpr(p, stmts, cn)
			return false

		case *ast.AssignStmt:
			if ts, containsError := getAstExprTypes(p, x.Rhs); containsError {
				if lhs, needReplace := normolizeAssignStmtTypes(p, x.Lhs, ts, fnScope.getNextErrorName()); needReplace {
					stmts[cn] = &replceCallExprStmt{
						nodeAfterInsertReturn: cn,
						lhs:                   lhs,
						rhs:                   x.Rhs,
						fnScope:               fnScope,
					}
				}
			}

		case *ast.ExprStmt:
			rhs := []ast.Expr{x.X}
			if ts, containsError := getAstExprTypes(p, []ast.Expr{x.X}); containsError {
				if lhs, needReplace := normolizeAssignStmtTypes(p, []ast.Expr{}, ts, fnScope.getNextErrorName()); needReplace {
					stmts[cn] = &replceCallExprStmt{
						nodeAfterInsertReturn: cn,
						lhs:                   lhs,
						rhs:                   rhs,
						fnScope:               fnScope,
					}
				}
			}
		}

		return true
	})
}
