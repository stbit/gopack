package syncerr

import (
	"go/ast"
	"strconv"

	"github.com/stbit/gopack/pkg/manager/pkginfo"
)

type functionScope struct {
	node    ast.Node
	errId   int
	results *ast.FieldList
}

func newFunctionScope(p *pkginfo.FileInfo, n ast.Node) *functionScope {
	fn := &functionScope{node: n}
	fn.results = fn.getResults()

	return fn
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
	r := f.results

	if r == nil {
		return false
	}

	if len(r.List) == 0 {
		return false
	}

	k := r.List[len(r.List)-1]

	if x, ok := k.Type.(*ast.Ident); ok && x.Name == "error" {
		return true
	}

	return false
}
