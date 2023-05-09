package parser

import "go/ast"

type funcWrapper struct {
	node *ast.Node
}

func newFuncWrapper(n *ast.Node) *funcWrapper {
	return &funcWrapper{node: n}
}

func (f *funcWrapper) getParams() *ast.FieldList {
	switch x := (*f.node).(type) {
	case *ast.FuncLit:
		return x.Type.Params
	case *ast.FuncDecl:
		return x.Type.Params
	}

	panic("not func declaration")
}

func (f *funcWrapper) getResults() *ast.FieldList {
	switch x := (*f.node).(type) {
	case *ast.FuncLit:
		return x.Type.Results
	case *ast.FuncDecl:
		return x.Type.Results
	}

	panic("not func declaration")
}

func (f *funcWrapper) getName() string {
	switch x := (*f.node).(type) {
	case *ast.FuncLit:
		return "anonimus"
	case *ast.FuncDecl:
		return x.Name.Name
	}

	panic("not func declaration")
}

func (f *funcWrapper) hasErrorResults() bool {
	r := f.getResults()

	if r == nil {
		return false
	}

	k := r.List[len(r.List)-1]
	l, ok := (k.Type).(*ast.Ident)
	if !ok {
		return false
	}

	return l.Name == "error"
}
