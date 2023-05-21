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
	fe := newFileInfoExtende(f)

	ast.Inspect(f.File, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncDecl:
			findReplacementExpr(fe, fe.stmts, n)
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

		if stmt, ok := fe.stmts[cn]; ok {
			stmt.replace(c)
			return false
		}

		return true
	}, nil)

	ast.Inspect(f.File, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			if x.Tok == token.IMPORT {
				if len(fe.zeroVariables) > 0 {
					iSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: strconv.Quote("reflect")}}
					x.Specs = append(x.Specs, iSpec)

					return false
				}
			}
		}

		return true
	})

	specs := fe.getZeroVariablesDecls()

	if len(specs) > 0 {
		f.File.Decls = append(f.File.Decls, &ast.GenDecl{
			Tok:   token.VAR,
			Specs: specs,
		})
	}
}

func findReplacementExpr(f *FileInfoExtended, stmts map[ast.Node]replceStmt, n ast.Node) {
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
					stmts[cn] = newReplceCallExprStmt(f, fnScope, cn, x.Lhs, x.Rhs)
				}
			}
		}

		return true
	})
}
