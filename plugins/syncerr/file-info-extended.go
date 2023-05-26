package syncerr

import (
	"go/ast"
	"strings"

	"github.com/stbit/gopack/pkg/manager/pkginfo"
)

type zeroValue struct {
	variable string
	typeVar  string
}

type fileInfoExtended struct {
	*pkginfo.FileContext
	stmts         map[ast.Node]replceStmt
	zeroVariables []zeroValue
}

func newFileInfoExtende(f *pkginfo.FileContext) *fileInfoExtended {
	return &fileInfoExtended{
		FileContext: f,
		stmts:       make(map[ast.Node]replceStmt),
	}
}

func (f *fileInfoExtended) getZeroVariablesDecls() []ast.Spec {
	var specs []ast.Spec = make([]ast.Spec, len(f.zeroVariables))
	for i, v := range f.zeroVariables {
		var t ast.Expr = ast.NewIdent(v.typeVar)

		if strings.Contains(v.typeVar, ".") {
			splits := strings.Split(v.typeVar, ".")
			t = &ast.SelectorExpr{
				Sel: ast.NewIdent(splits[1]),
				X:   ast.NewIdent(splits[0]),
			}
		}

		specs[i] = &ast.ValueSpec{
			Names: []*ast.Ident{{Name: v.variable}},
			Type:  t,
		}
	}

	return specs
}
