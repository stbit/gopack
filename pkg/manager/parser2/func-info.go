package parser2

import (
	"fmt"
	"go/ast"
	"go/types"
)

func getFuncResults(p *SourcePackage, c *ast.CallExpr) *types.Tuple {
	switch s := c.Fun.(type) {
	case *ast.Ident:
		if !p.pkg.TypesInfo.Types[s].IsType() {
			if use, ok := p.pkg.TypesInfo.Uses[s]; ok {
				fmt.Println(s)
				return use.Type().(*types.Signature).Results()
			}
		}
	case *ast.SelectorExpr:
		if s.Sel != nil && !p.pkg.TypesInfo.Types[s.Sel].IsType() {
			if use, ok := p.pkg.TypesInfo.Uses[s.Sel]; ok {
				fmt.Println(s)
				return use.Type().(*types.Signature).Results()
			}
		}
	}

	return nil
}

func hasFuncResultsError(p *SourcePackage, c *ast.CallExpr) (*types.Tuple, bool) {
	r := getFuncResults(p, c)

	if r.Len() > 0 {
		l := r.At(r.Len() - 1)

		return r, l.Type().String() == "error"
	}

	return r, false
}

func normolizeAssignStmt(p *SourcePackage, exps []ast.Expr, fresults *types.Tuple, errName string) ([]ast.Expr, bool) {
	r := make([]ast.Expr, fresults.Len())
	needReplace := false

	for i := 0; i < fresults.Len(); i++ {
		n := "_"

		if len(exps)-1 >= i {
			ind := exps[i].(*ast.Ident)

			if ind == nil {
				panic(fmt.Errorf("expression can only ast.Ident"))
			}

			n = ind.Name
		}

		if fresults.Len()-1 == i && n == "_" {
			n = errName
			needReplace = true
		}

		r[i] = &ast.Ident{
			Name: n,
		}
	}

	return r, needReplace
}
