package manager

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

// var re = regexp.MustCompile(`(\s+)(.*?)?_\s+?:?=\s+?(.*?)\r\n`)

type ParseFile struct {
	Path string
}

func newParseFile(path string) *ParseFile {
	return &ParseFile{Path: path}
}

func (f *ParseFile) parse(dist string) error {
	content, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}

	fset, file := parentFunc(string(content))

	of, err := os.OpenFile(dist, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}

	defer of.Close()

	printer.Fprint(of, fset, file)

	return err
}

func parentFunc(content string) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "demo", content, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		funcCall, ok := n.(*ast.CallExpr)
		if ok {
			mthd, ok := funcCall.Fun.(*ast.SelectorExpr)
			if ok {
				_, ok := mthd.X.(*ast.Ident)
				if ok {
					// fmt.Printf(in.Name)
				}
			}
		}

		return true
	})

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.FuncDecl:
			c.Replace(replaceErrors(&n, x.Name.Name == "main"))
		}

		return true
	})

	// printer.Fprint(os.Stdout, fset, file)

	return fset, file
}

func getReturnErrorStm(nameErr string, isPanic bool) ast.Stmt {
	if isPanic {
		return &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.Ident{Name: "panic"},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: nameErr,
					},
				},
			},
		}
	} else {
		return &ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.Ident{Name: nameErr},
			},
		}
	}
}

func replaceErrors(n *ast.Node, isPanic bool) ast.Node {
	errInc := 0
	return astutil.Apply(*n, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
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
							getReturnErrorStm(v.Name, isPanic),
						},
					},
				})
			}
		}

		return true
	})
}
