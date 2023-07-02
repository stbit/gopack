package syncerr

import (
	"fmt"
	"go/token"
	"strconv"

	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

var zeroVariableId = 0

type replceCallExprStmt struct {
	nodeAfterInsertReturn dst.Node
	fileInfo              *fileInfoExtended
	fnScope               *functionScope
	lhs                   []dst.Expr
	rhs                   []dst.Expr
}

func newReplceCallExprStmt(fe *fileInfoExtended, f *functionScope, n dst.Node, lhs []dst.Expr, rhs []dst.Expr) *replceCallExprStmt {
	return &replceCallExprStmt{
		nodeAfterInsertReturn: n,
		lhs:                   lhs,
		rhs:                   rhs,
		fileInfo:              fe,
		fnScope:               f,
	}
}

func (s *replceCallExprStmt) replace(c *dstutil.Cursor) {
	if !s.fnScope.hasErrorResults() {
		panic(fmt.Errorf("func %s not return error", s.fnScope.getName()))
	}

	errIdent := s.lhs[len(s.lhs)-1].(*dst.Ident)
	errIdent.Name = s.fnScope.getNextErrorName()
	ts := s.fnScope.getResults()
	tslen := len(ts.List)
	results := make([]dst.Expr, len(ts.List))

	for i, t := range ts.List {
		results[i] = &dst.Ident{Name: s.getDefaultValue(t, errIdent.Name, tslen-1 == i)}
	}

	c.Replace(&dst.AssignStmt{
		Lhs: s.lhs,
		Tok: token.DEFINE,
		Rhs: s.rhs,
	})

	nameErr := s.lhs[len(s.lhs)-1].(*dst.Ident).Name

	c.InsertAfter(&dst.IfStmt{
		Cond: &dst.BinaryExpr{
			// err
			X: &dst.Ident{Name: nameErr},
			// !=
			Op: token.NEQ,
			// nil
			Y: &dst.Ident{Name: "nil"},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: results,
				},
			},
		},
	})
}

func (r *replceCallExprStmt) getZeroValue(name string) string {
	var zv zeroValue

	for _, v := range r.fileInfo.zeroVariables {
		if v.typeVar == name {
			zv = v
			break
		}
	}

	if zv.typeVar == "" {
		zeroVariableId++
		zv = zeroValue{"zdv_" + strconv.Itoa(zeroVariableId), name}
		r.fileInfo.zeroVariables = append(r.fileInfo.zeroVariables, zv)
	}

	return zv.variable
}

func (r *replceCallExprStmt) getDefaultValue(f *dst.Field, errName string, isLast bool) string {
	switch x := f.Type.(type) {
	case *dst.StarExpr, *dst.ArrayType, *dst.FuncType:
		return "nil"

	case *dst.Ident:
		if isLast && x.Name == "error" {
			return errName
		} else {
			switch x.Name {
			case "bool":
				return "false"
			}

			return r.getZeroValue(x.Name)
		}

	case *dst.SelectorExpr:
		if i, ok := x.X.(*dst.Ident); ok {
			return r.getZeroValue(i.Name + "." + x.Sel.Name)
		}
	}

	return "nil"
}
