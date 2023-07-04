package errstack

import (
	"go/token"
	"strconv"

	"github.com/dave/dst"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
)

var variableId = 0

type fileInfoExtended struct {
	*pkginfo.FileContext
	hasWrap bool
	varName string
}

func newFileInfoExtende(f *pkginfo.FileContext) *fileInfoExtended {
	variableId++
	return &fileInfoExtended{
		FileContext: f,
		varName:     "errstack_empty_" + strconv.Itoa(variableId),
	}
}

func (f *fileInfoExtended) SetWrap() {
	f.hasWrap = true
}

func (f *fileInfoExtended) getVariablesDecls() []dst.Spec {
	if !f.hasWrap {
		return []dst.Spec{}
	}

	specs := []dst.Spec{
		&dst.ValueSpec{
			Names: []*dst.Ident{dst.NewIdent(f.varName)},
			Type:  dst.NewIdent("string"),
			Values: []dst.Expr{
				&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote("")},
			},
		},
	}

	return specs
}
