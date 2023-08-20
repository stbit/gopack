package errstack

import (
	"go/printer"
	"go/token"
	"os"

	"github.com/dave/dst"
	"github.com/stbit/gopack/pkg/manager/hooks"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
	"github.com/stbit/gopack/plugins"
)

var pluginName string = "errstack"

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
			findFuncExpr(fe, n)
			return false
		}

		return true
	})

	specs := fe.getVariablesDecls()

	if len(specs) > 0 {
		f.File.Decls = append(f.File.Decls, &dst.GenDecl{
			Tok:   token.VAR,
			Specs: specs,
		})

		f.AddImport("github.com/rotisserie/eris")
	}
}

func findFuncExpr(f *fileInfoExtended, n dst.Node) {
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
			findFuncExpr(f, cn)
			return false

		case *dst.ReturnStmt:
			if x.Results != nil {
				last := x.Results[len(x.Results)-1]

				if fnScope.hasErrorResults() && canAddWrap(last) && len(fnScope.getResults().List) == len(fnScope.getParams().List) {
					f.SetWrap()
					x.Results[len(x.Results)-1] = &dst.CallExpr{
						Args: []dst.Expr{
							last,
							dst.NewIdent(f.varName),
						},
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("eris"),
							Sel: dst.NewIdent("Wrap"),
						},
					}
				}
			}
		}

		return true
	})
}

func canAddWrap(exp dst.Expr) bool {
	switch x := exp.(type) {
	case *dst.Ident:
		return x.Name != "nil"
	case *dst.CallExpr:
		if c, ok := x.Fun.(*dst.SelectorExpr); ok {
			if i, ok := c.X.(*dst.Ident); ok {
				return i.Name != "eris"
			}
		}

	}

	return true
}
