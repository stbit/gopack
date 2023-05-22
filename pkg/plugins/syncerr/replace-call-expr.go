package syncerr

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

var (
	zeroValue      = []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune", "float32", "float64", "complex64"}
	zeroVariableId = 0
)

type replceCallExprStmt struct {
	nodeAfterInsertReturn ast.Node
	fileInfo              *FileInfoExtended
	fnScope               *functionScope
	lhs                   []ast.Expr
	rhs                   []ast.Expr
}

func newReplceCallExprStmt(fe *FileInfoExtended, f *functionScope, n ast.Node, lhs []ast.Expr, rhs []ast.Expr) *replceCallExprStmt {
	return &replceCallExprStmt{
		nodeAfterInsertReturn: n,
		lhs:                   lhs,
		rhs:                   rhs,
		fileInfo:              fe,
		fnScope:               f,
	}
}

func (s *replceCallExprStmt) replace(c *astutil.Cursor) {
	if !s.fnScope.hasErrorResults() {
		panic(fmt.Errorf("func %s not return error", s.fnScope.getName()))
	}

	errIdent := s.lhs[len(s.lhs)-1].(*ast.Ident)
	errIdent.Name = s.fnScope.getNextErrorName()
	ts := s.fnScope.getResults()
	tslen := len(ts.List)
	results := make([]ast.Expr, len(ts.List))

	for i, t := range ts.List {
		results[i] = &ast.Ident{Name: s.getDefaultValue(t, errIdent.Name, tslen-1 == i)}
	}

	c.Replace(&ast.AssignStmt{
		Lhs: s.lhs,
		Tok: token.DEFINE,
		Rhs: s.rhs,
	})

	nameErr := s.lhs[len(s.lhs)-1].(*ast.Ident).Name

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

func (r *replceCallExprStmt) getZeroValue(name string) string {
	var zv ZeroValue

	for _, v := range r.fileInfo.zeroVariables {
		if v.typeVar == name {
			zv = v
			break
		}
	}

	if zv.typeVar == "" {
		zeroVariableId++
		zv = ZeroValue{"zdv_" + strconv.Itoa(zeroVariableId), name}
		r.fileInfo.zeroVariables = append(r.fileInfo.zeroVariables, zv)
	}

	return zv.variable
}

func (r *replceCallExprStmt) getDefaultValue(f *ast.Field, errName string, isLast bool) string {
	switch x := f.Type.(type) {
	case *ast.StarExpr, *ast.ArrayType, *ast.FuncType:
		return "nil"

	case *ast.Ident:
		if isLast && x.Name == "error" {
			return errName
		} else {
			for _, v := range zeroValue {
				if v == x.Name {
					return "0"
				}
			}

			switch x.Name {
			case "string":
				return "\"\""
			case "bool":
				return "false"
			}

			return r.getZeroValue(x.Name)
		}

	case *ast.SelectorExpr:
		if i, ok := x.X.(*ast.Ident); ok {
			return r.getZeroValue(i.Name + "." + x.Sel.Name)
		}
	}

	return "nil"
}
