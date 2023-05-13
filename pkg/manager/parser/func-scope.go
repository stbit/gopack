package parser

import (
	"go/ast"
	"strconv"
)

type functionScope struct {
	node  ast.Node
	errId int
}

func newFunctionScope(p *SourcePackage, n ast.Node) *functionScope {
	return &functionScope{node: n}
}

func (f *functionScope) getResults() *ast.FieldList {
	switch x := f.node.(type) {
	case *ast.FuncLit:
		return x.Type.Results
	case *ast.FuncDecl:
		return x.Type.Results
	}

	panic("not func declaration")
}

func (f *functionScope) getName() string {
	switch x := f.node.(type) {
	case *ast.FuncLit:
		return "anonimus"
	case *ast.FuncDecl:
		return x.Name.Name
	}

	panic("not func declaration")
}

func (f *functionScope) getNextErrorName() string {
	f.errId++
	return "err_" + strconv.Itoa(f.errId)
}

func (f *functionScope) hasErrorResults() bool {
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
