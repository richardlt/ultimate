package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
	"github.com/richardlt/ultimate"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing path of files to parse")
	}
	for _, p := range os.Args[1:] {
		log.Printf("parsing file at '%s'", p)
		if err := parse(p); err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func parse(path string) error {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "cannot read file at '%s'", path)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return errors.Wrapf(err, "cannot parse file at '%s'", path)
	}

	structs := make(map[string]*ast.StructType)
	comments := make(map[string]string)
	var comment string
	ast.Inspect(file, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.Comment:
			if strings.HasPrefix(v.Text, "//ultimate") {
				comment = v.Text
			}
		case *ast.TypeSpec:
			if s, ok := v.Type.(*ast.StructType); ok {
				structs[v.Name.Name] = s
				comments[v.Name.Name] = comment
				comment = ""
				return false
			}
		default:
			if v != nil && comment != "" {
				comment = ""
			}
		}
		return true
	})

	var crs []*ultimate.Criteria

	for n, s := range structs {
		if comments[n] != "" {
			cr, err := extractCriteria(src, n, comments[n], s)
			if err != nil {
				return errors.Wrapf(err, "cannot generate criteria for struct '%s'", n)
			}
			crs = append(crs, cr)
		}
	}

	for _, cr := range crs {
		src, err := cr.Generate()
		if err != nil {
			return err
		}
		fmt.Println(src)
	}

	return nil
}

func extractCriteria(src []byte, name, comment string, s *ast.StructType) (*ultimate.Criteria, error) {
	split := strings.Split(comment, ":")
	if len(split) != 2 {
		return nil, errors.Errorf("invalid ultimate comment value for struct %s", name)
	}

	ct, err := ultimate.ParseCriteriaType(split[1])
	if err != nil {
		return nil, err
	}

	c := ultimate.NewCriteria(name, ct)

	for _, field := range s.Fields.List {
		if field.Tag == nil {
			continue
		}

		name := field.Names[0].Name

		ts, err := structtag.Parse(strings.Replace(field.Tag.Value, "`", "", -1))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse tags for field '%s'", name)
		}

		tag, err := ts.Get("ultimate")
		if err == nil {
			for _, o := range tag.Options {
				if o == "criteria" {
					start, end := field.Type.Pos()-1, field.Type.End()-1
					fieldType, err := ultimate.ParseCriteriaFieldType(string(src[start:end]))
					if err != nil {
						return nil, err
					}
					c.AddField(tag.Name, fieldType)
				}
			}
		}
	}

	return &c, nil
}
