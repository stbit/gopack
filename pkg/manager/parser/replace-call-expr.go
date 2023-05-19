package parser

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type replceCallExprStmt struct {
	nodeAfterInsertReturn ast.Node
	fnScope               *functionScope
	lhs                   []ast.Expr
	rhs                   []ast.Expr
}

func (s *replceCallExprStmt) replace(c *astutil.Cursor) {
	if !s.fnScope.hasErrorResults() {
		panic(fmt.Errorf("func %s not return error", s.fnScope.getName()))
	}

	nameErr := s.lhs[len(s.lhs)-1].(*ast.Ident).Name
	c.Replace(&ast.AssignStmt{
		Lhs: s.lhs,
		Tok: token.DEFINE,
		Rhs: s.rhs,
	})

	ts := s.fnScope.returnTypes
	tslen := len(ts)
	results := make([]ast.Expr, len(ts))

	for i, t := range ts {
		if t != nil {
			if tslen-1 == i && t.String() == "error" {
				results[i] = &ast.Ident{Name: nameErr}
			} else {
				results[i] = &ast.Ident{Name: getDefaultValue(s.fnScope.getTypeName(t))}
			}
		} else {
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
