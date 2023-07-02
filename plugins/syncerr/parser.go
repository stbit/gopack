package syncerr

import (
	"go/printer"
	"go/token"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
	"github.com/stbit/gopack/pkg/manager/hooks"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
	"github.com/stbit/gopack/plugins"
)

var pluginName = "syncerr"

type replceStmt interface {
	replace(c *dstutil.Cursor)
}

func New() plugins.PluginRegister {
	return func(m *plugins.ManagerContext) error {
		m.AddHookParseFile(pluginName, hooks.HOOK_PARSE_FILE, func(f *pkginfo.FileContext) error {
			parseFile(f)
			return nil
		})

		return nil
	}
}

func parseFile(f *pkginfo.FileContext) {
	fe := newFileInfoExtende(f)

	dst.Inspect(f.File, func(n dst.Node) bool {
		switch n.(type) {
		case *dst.FuncDecl:
			findReplacementExpr(fe, fe.stmts, n)
			return false
		}

		return true
	})

	dstutil.Apply(f.File, func(c *dstutil.Cursor) bool {
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

	specs := fe.getZeroVariablesDecls()

	if len(specs) > 0 {
		f.File.Decls = append(f.File.Decls, &dst.GenDecl{
			Tok:   token.VAR,
			Specs: specs,
		})
	}
}

func findReplacementExpr(f *fileInfoExtended, stmts map[dst.Node]replceStmt, n dst.Node) {
	fnScope := newFunctionScope(f, n)
	dst.Inspect(n, func(cn dst.Node) bool {
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
		case *dst.FuncDecl, *dst.FuncLit:
			findReplacementExpr(f, stmts, cn)
			return false

		case *dst.AssignStmt:
			switch i := x.Lhs[len(x.Lhs)-1].(type) {
			case *dst.Ident:
				if i.Name == "_" {
					stmts[cn] = newReplceCallExprStmt(f, fnScope, cn, x.Lhs, x.Rhs)
				}
			}
		}

		return true
	})
}
