package jsontag

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"github.com/fatih/structtag"
	"github.com/stbit/gopack/pkg/manager/hooks"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
	"github.com/stbit/gopack/plugins"
)

var pluginName = "jsontag"

type Transformer = func(tag *Tag) (name string, options []string)

type Tag struct {
	FieldName  string
	JsonName   string
	Options    []string
	FilePath   string
	StructName string
}

func New(transform Transformer) plugins.PluginRegister {
	return func(m *plugins.ManagerContext) error {
		m.AddHookParseFile(pluginName, hooks.HOOK_PARSE_FILE, func(f *pkginfo.FileContext) error {
			recurciveReplaceStructTags(f, f.File, transform)

			return nil
		})

		return nil
	}
}

func recurciveReplaceStructTags(f *pkginfo.FileContext, cn ast.Node, transform Transformer) {
	ast.Inspect(cn, func(n ast.Node) bool {
		var (
			structName string
			t          *ast.StructType
		)

		switch s := n.(type) {
		case *ast.TypeSpec:
			if k, ok := s.Type.(*ast.StructType); ok {
				structName = s.Name.Name
				t = k
			}
		case *ast.Field:
			if k, ok := s.Type.(*ast.StructType); ok {
				structName = s.Names[0].Name
				t = k
			}
		case *ast.StructType:
			t = s
		}

		if t != nil {
			tag := &Tag{FilePath: f.GetSourcePath(), StructName: structName}
			for _, v := range t.Fields.List {
				if len(v.Names) > 0 {
					recurciveReplaceStructTags(f, v, transform)

					tag.FieldName = v.Names[0].Name
					err := replaceFieldJsonName(v, tag, transform)
					if err != nil {
						tagName := ""
						if v.Tag != nil {
							tagName = v.Tag.Value
						}

						f.AddError(fmt.Errorf("file(%s) struct(%s) field(%s) tag(%s): %v", tag.FilePath, tag.StructName, tag.FieldName, tagName, err))
					}
				}
			}

			return false
		}

		return true
	})
}

func replaceFieldJsonName(n *ast.Field, t *Tag, transform Transformer) error {
	tagsStr := ""

	if n.Tag != nil {
		tagsStr = strings.ReplaceAll(n.Tag.Value, "`", "")
	}

	tags, err := structtag.Parse(tagsStr)
	if err != nil {
		return err
	}

	jsonTag, err := tags.Get("json")
	if err != nil {
		t.Options = make([]string, 0)
		nn, no := transform(t)
		jsonTag = &structtag.Tag{
			Key:     "json",
			Name:    nn,
			Options: no,
		}
	} else {
		t.JsonName = jsonTag.Name
		t.Options = jsonTag.Options
		nn, no := transform(t)
		jsonTag.Name = nn
		jsonTag.Options = no
	}

	tags.Set(jsonTag)
	sort.Sort(tags)
	jsonVal := fmt.Sprintf("`%s`", tags.String())

	if n.Tag != nil {
		n.Tag.Value = jsonVal
	} else {
		n.Tag = &ast.BasicLit{Kind: token.STRING, Value: jsonVal}
	}

	return nil
}
