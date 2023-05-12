package parser

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type replceCallExprStmt struct {
	parentNode ast.Node
	callExpr   *ast.CallExpr
	fw         *funcWrapper
	lhs        []ast.Expr
}

func (s *replceCallExprStmt) replace(c *astutil.Cursor) {
	if !s.fw.hasErrorResults() {
		panic(fmt.Errorf("func %s not return error", s.fw.getName()))
	}

	nameErr := s.lhs[len(s.lhs)-1].(*ast.Ident).Name
	c.Replace(&ast.AssignStmt{
		Lhs: s.lhs,
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			s.callExpr,
		},
	})

	fwResults := s.fw.getResults().List
	results := make([]ast.Expr, len(fwResults))

	for i, v := range fwResults {
		switch x := v.Type.(type) {
		case *ast.StarExpr:
			results[i] = &ast.Ident{Name: "nil"}
		case *ast.Ident:
			if x.Name == "error" {
				results[i] = &ast.Ident{Name: nameErr}
			} else {
				results[i] = &ast.Ident{Name: getDefaultValue(x.Name)}
			}
		default:
			results[i] = &ast.Ident{Name: "nil"}
		}
	}

	c.InsertAfter(&ast.IfStmt{
		Cond: &ast.BinaryExpr{
			// err
			X: &ast.Ident{Name: nameErr},
			// !=
			Op: token.NEQ,
			// nil
			Y: &ast.Ident{Name: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: results,
				},
			},
		},
	})
}
