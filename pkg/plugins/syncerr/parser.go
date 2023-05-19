package syncerr

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"strconv"

	"github.com/stbit/gopack/pkg/manager/pkginfo"
	"golang.org/x/tools/go/ast/astutil"
)

type replceStmt interface {
	replace(c *astutil.Cursor)
}

func ParseFile(f *pkginfo.FileInfo) {
	stmts := make(map[ast.Node]replceStmt)

	ast.Inspect(f.File, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncDecl:
			findReplacementExpr(f, stmts, n)
			return false
		}

		return true
	})

	astutil.Apply(f.File, func(c *astutil.Cursor) bool {
		cn := c.Node()

		defer func() {
			if r := recover(); r != nil {
				if err := printer.Fprint(os.Stdout, f.Fset, cn); err != nil {
					panic(err)
				}

				panic(r)
			}
		}()

		if stmt, ok := stmts[cn]; ok {
			stmt.replace(c)
			return false
		}

		return true
	}, nil)

	ast.Inspect(f.File, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			if x.Tok == token.IMPORT {
				for _, v := range stmts {
					if ex, ok := v.(*replceCallExprStmt); ok && len(ex.zeroNameVars) > 0 {
						// Add the new import
						iSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: strconv.Quote("reflect")}}
						x.Specs = append(x.Specs, iSpec)
						return false
					}
				}
			}

			return false
		}

		return true
	})
}

func findReplacementExpr(f *pkginfo.FileInfo, stmts map[ast.Node]replceStmt, n ast.Node) {
	fnScope := newFunctionScope(f, n)
	ast.Inspect(n, func(cn ast.Node) bool {
		defer func() {
			if r := recover(); r != nil {
				if err := printer.Fprint(os.Stdout, f.Fset, cn); err != nil {
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
			findReplacementExpr(f, stmts, cn)
			return false

		case *ast.AssignStmt:
			switch i := x.Lhs[len(x.Lhs)-1].(type) {
			case *ast.Ident:
				if i.Name == "_" {
					stmts[cn] = newReplceCallExprStmt(fnScope, cn, x.Lhs, x.Rhs)
				}
			}
		}

		return true
	})
}
