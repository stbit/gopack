package parser

import (
	"go/ast"
	"go/types"
	"strconv"
	"strings"
)

type functionScope struct {
	node        ast.Node
	errId       int
	returnTypes []types.Type
}

func newFunctionScope(p *SourcePackage, n ast.Node) *functionScope {
	fn := &functionScope{node: n, returnTypes: []types.Type{}}
	r := fn.getResults()

	if r != nil && r.List != nil {
		for _, v := range r.List {
			appened := false

			switch x := v.Type.(type) {
			case *ast.Ident:
				if use, ok := p.pkg.TypesInfo.Uses[x]; ok {
					if use.Type().Underlying().String() != "interface{}" {
						appened = true
						fn.returnTypes = append(fn.returnTypes, use.Type())
					}
				}
			}

			if !appened {
				fn.returnTypes = append(fn.returnTypes, nil)
			}
		}
	}

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
	r := f.returnTypes

	if len(r) == 0 {
		return false
	}

	k := r[len(r)-1]

	if k == nil {
		return false
	}

	return k.String() == "error"
}

func (f *functionScope) getTypeName(t types.Type) string {
	n := t.String()
	d := strings.LastIndex(n, ".")

	if d != -1 {
		return string(n[d+1:])
	} else {
		return n
	}
}
