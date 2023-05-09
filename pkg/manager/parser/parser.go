package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

type ParseFile struct {
	Path string
}

func NewParseFile(path string) *ParseFile {
	return &ParseFile{Path: path}
}

func (f *ParseFile) Parse(dist string) error {
	content, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}

	fset, file := parentFunc(f.Path, string(content))

	of, err := os.OpenFile(dist, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}

	defer of.Close()

	printer.Fprint(of, fset, file)

	return err
}

func parentFunc(filePath string, content string) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch n.(type) {
		case *ast.FuncDecl:
			c.Replace(replaceErrors(&n))
		}

		return true
	})

	return fset, file
}

func getReturnErrorStm(fw *funcWrapper, nameErr string) ast.Stmt {
	fmt.Println(fw.getName())

	if !fw.hasErrorResults() {
		panic(fmt.Errorf("func %s not return error", fw.getName()))
	}

	results := make([]ast.Expr, 0)

	for _, v := range fw.getResults().List {
		switch x := v.Type.(type) {
		case *ast.StarExpr:
			results = append(results, &ast.Ident{Name: "nil"})
		case *ast.Ident:
			if x.Name == "error" {
				results = append(results, &ast.Ident{Name: nameErr})
			} else {
				results = append(results, &ast.Ident{Name: getDefaultValue(x.Name)})
			}
		default:
			results = append(results, &ast.Ident{Name: "nil"})
		}
	}

	return &ast.ReturnStmt{
		Results: results,
	}
}

func replaceErrors(n *ast.Node) ast.Node {
	funcWrapper := newFuncWrapper(n)
	errInc := 0
	return astutil.Apply(*n, func(c *astutil.Cursor) bool {
		cn := c.Node()

		if *n == cn {
			return true
		}

		switch x := cn.(type) {
		case *ast.CallExpr:
			id, ok := x.Fun.(*ast.Ident)

			if ok {
				if id.Name == "funcWithOneError" {
					c.Replace(&ast.UnaryExpr{
						Op: token.NOT,
						X:  x,
					})
				}
			}
		case *ast.FuncLit:
			c.Replace(replaceErrors(&cn))
			return false
		case *ast.AssignStmt:
			l := x.Lhs[len(x.Lhs)-1]
			v, ok := l.(*ast.Ident)

			if ok && v.Name == "_" {

				errInc++
				v.Name = "err_" + strconv.Itoa(errInc)

				c.InsertAfter(&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						// err
						X: &ast.Ident{Name: v.Name},
						// !=
						Op: token.NEQ,
						// nil
						Y: &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							getReturnErrorStm(funcWrapper, v.Name),
						},
					},
				})
			}

		// exit if nested func declaration
		case *ast.FuncDecl:
			return false
		}

		return true
	}, nil)
}
