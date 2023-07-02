package syncerr

import (
	"strings"

	"github.com/dave/dst"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
)

type zeroValue struct {
	variable string
	typeVar  string
}

type fileInfoExtended struct {
	*pkginfo.FileContext
	stmts         map[dst.Node]replceStmt
	zeroVariables []zeroValue
}

func newFileInfoExtende(f *pkginfo.FileContext) *fileInfoExtended {
	return &fileInfoExtended{
		FileContext: f,
		stmts:       make(map[dst.Node]replceStmt),
	}
}

func (f *fileInfoExtended) getZeroVariablesDecls() []dst.Spec {
	var specs []dst.Spec = make([]dst.Spec, len(f.zeroVariables))
	for i, v := range f.zeroVariables {
		var t dst.Expr = dst.NewIdent(v.typeVar)

		if strings.Contains(v.typeVar, ".") {
			splits := strings.Split(v.typeVar, ".")
			t = &dst.SelectorExpr{
				Sel: dst.NewIdent(splits[1]),
				X:   dst.NewIdent(splits[0]),
			}
		}

		specs[i] = &dst.ValueSpec{
			Names: []*dst.Ident{{Name: v.variable}},
			Type:  t,
		}
	}

	return specs
}
