package errstack

import (
	"go/token"
	"strconv"

	"github.com/dave/dst"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
)

type fileInfoExtended struct {
	*pkginfo.FileContext
	hasWrap bool
}

func newFileInfoExtende(f *pkginfo.FileContext) *fileInfoExtended {
	return &fileInfoExtended{
		FileContext: f,
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
			Names: []*dst.Ident{dst.NewIdent("errstack_empty")},
			Type:  dst.NewIdent("string"),
			Values: []dst.Expr{
				&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote("")},
			},
		},
	}

	return specs
}
