package parser

import (
	"fmt"
	"go/ast"
	"go/types"
)

func resolveExprTypes(p *SourcePackage, id *ast.Ident) ([]types.Type, bool, error) {
	hasError := false

	if !p.pkg.TypesInfo.Types[id].IsType() {
		if use, ok := p.pkg.TypesInfo.Uses[id]; ok {
			switch x := use.Type().(type) {
			case *types.Signature:
				rf := x.Results()
				r := make([]types.Type, rf.Len())
				for i := 0; i < rf.Len(); i++ {
					r[i] = rf.At(i).Type()

					if !hasError {
						hasError = r[i].String() == "error"
					}
				}

				return r, hasError, nil
			default:
				return []types.Type{x}, x.String() == "error", nil
			}
		}
	}

	return nil, false, fmt.Errorf("cant find type")
}

func getAstExprTypes(p *SourcePackage, exprs []ast.Expr) ([][]types.Type, bool) {
	var ind *ast.Ident
	result := make([][]types.Type, len(exprs))
	rhasError := false

	for i, e := range exprs {
		ind = nil

		switch x := e.(type) {
		case *ast.SelectorExpr:
			ind = x.Sel
		case *ast.CallExpr:
			switch c := x.Fun.(type) {
			case *ast.SelectorExpr:
				ind = c.Sel
			case *ast.Ident:
				ind = c
			}
		case *ast.Ident:
			ind = x
		default:
			result[i] = []types.Type{&types.Basic{}}
		}

		if ind != nil {
			r, hasError, err := resolveExprTypes(p, ind)
			if err != nil {
				panic(err)
			}

			if hasError {
				rhasError = true
			}

			result[i] = r
		}
	}

	return result, rhasError
}

func normolizeAssignStmtTypes(p *SourcePackage, exps []ast.Expr, rhstypes [][]types.Type, errName string) ([]ast.Expr, bool) {
	var ftypes []types.Type

	if len(rhstypes) == 1 {
		ftypes = rhstypes[0]
	} else {
		ftypes = make([]types.Type, len(rhstypes))

		for i, c := range rhstypes {
			if len(c) > 1 {
				panic("multiple-value in single-value context")
			}

			ftypes[i] = c[0]
		}
	}

	le := len(exps)
	lr := len(ftypes)
	r := make([]ast.Expr, lr)
	needReplace := false

	if le > 0 && le != lr {
		panic("left exprs count not equal right")
	}

	for i := 0; i < lr; i++ {
		if i < le {
			r[i] = exps[i]
		} else {
			r[i] = &ast.Ident{
				Name: "_",
			}
		}

		if lr-1 == i {
			if n, ok := r[i].(*ast.Ident); ok && n.Name == "_" {
				n.Name = errName
				needReplace = true
			}
		}
	}

	return r, needReplace
}
