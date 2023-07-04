package errstack

import (
	"github.com/dave/dst"
)

type functionScope struct {
	node    dst.Node
	results *dst.FieldList
}

func newFunctionScope(p *fileInfoExtended, n dst.Node) *functionScope {
	fn := &functionScope{node: n}
	fn.results = fn.getResults()

	return fn
}

func (f *functionScope) getResults() *dst.FieldList {
	switch x := f.node.(type) {
	case *dst.FuncLit:
		return x.Type.Results
	case *dst.FuncDecl:
		return x.Type.Results
	}

	panic("not func declaration")
}

func (f *functionScope) getName() string {
	switch x := f.node.(type) {
	case *dst.FuncLit:
		return "anonimus"
	case *dst.FuncDecl:
		return x.Name.Name
	}

	panic("not func declaration")
}

func (f *functionScope) hasErrorResults() bool {
	r := f.results

	if r == nil {
		return false
	}

	if len(r.List) == 0 {
		return false
	}

	k := r.List[len(r.List)-1]

	if x, ok := k.Type.(*dst.Ident); ok && x.Name == "error" {
		return true
	}

	return false
}
