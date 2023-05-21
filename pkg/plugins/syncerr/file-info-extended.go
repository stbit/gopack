package syncerr

import (
	"go/ast"

	"github.com/stbit/gopack/pkg/manager/pkginfo"
)

type ZeroValue struct {
	variable string
	typeVar  string
	expr     string
}

type FileInfoExtended struct {
	*pkginfo.FileInfo
	stmts         map[ast.Node]replceStmt
	zeroVariables []ZeroValue
}

func newFileInfoExtende(f *pkginfo.FileInfo) *FileInfoExtended {
	return &FileInfoExtended{
		FileInfo: f,
		stmts:    make(map[ast.Node]replceStmt),
	}
}

func (f *FileInfoExtended) getZeroVariablesDecls() []ast.Spec {
	var specs []ast.Spec = make([]ast.Spec, len(f.zeroVariables))
	for i, v := range f.zeroVariables {
		specs[i] = &ast.ValueSpec{
			Names: []*ast.Ident{{Name: v.variable}},
			Values: []ast.Expr{&ast.TypeAssertExpr{
				Type: ast.NewIdent(v.typeVar),
				X: &ast.CallExpr{Fun: &ast.SelectorExpr{
					Sel: ast.NewIdent("Interface"),
					X: &ast.CallExpr{
						Args: []ast.Expr{&ast.CallExpr{Fun: &ast.SelectorExpr{
							Sel: ast.NewIdent("Elem"),
							X: &ast.CallExpr{
								Args: []ast.Expr{
									&ast.CallExpr{
										Args: []ast.Expr{ast.NewIdent("nil")},
										Fun:  &ast.ParenExpr{X: &ast.StarExpr{X: ast.NewIdent(v.typeVar)}},
									},
								},
								Fun: &ast.SelectorExpr{
									Sel: ast.NewIdent("TypeOf"),
									X:   ast.NewIdent("reflect"),
								},
							},
						}}},
						Fun: &ast.SelectorExpr{
							Sel: ast.NewIdent("Zero"),
							X:   ast.NewIdent("reflect"),
						},
					},
				}},
			}},
		}
	}

	return specs
}
