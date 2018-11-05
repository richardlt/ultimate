package ultimate

import (
	"bytes"
	"strings"

	"html/template"

	"github.com/pkg/errors"
)

type CriteriaType int

const (
	CriteriaSQL CriteriaType = iota
	CriteriaMongo
)

func ParseCriteriaType(t string) (CriteriaType, error) {
	switch t {
	case "sql":
		return CriteriaSQL, nil
	case "mongo":
		return CriteriaMongo, nil
	}
	return 0, errors.Errorf("invalid given criteria type value '%s'", t)
}

func NewCriteria(name string, t CriteriaType) Criteria {
	return Criteria{Name: name, Type: t}
}

type Criteria struct {
	Name   string
	Type   CriteriaType
	Fields []CriteriaField
}

func (c *Criteria) AddField(name string, fieldType CriteriaFieldType) {
	c.Fields = append(c.Fields, CriteriaField{Name: name, Type: fieldType})
}

func (c Criteria) Generate() (string, error) {
	content := `
package {{.Name | ToLower}}

func NewCriteria{{.Name | Title}}() *Criteria{{.Name | Title}} {
	return &Criteria{{.Name | Title}}{}
}

type Criteria{{.Name | Title}} struct { {{range .Fields }}
	{{.Name}}s {{.Type}}{{end}}
}
{{range .Fields}}
func (c *Criteria{{$.Name | Title}}) {{.Name | Title }}s(v {{.Type}}...) *Criteria{{$.Name | Title}}{
	c.{{.Name}}s = v
	return c
}
{{end}}
	`

	funcMap := template.FuncMap{
		"Title":   strings.Title,
		"ToLower": strings.ToLower,
	}

	tmpl, err := template.New("criteria").Funcs(funcMap).Parse(content)
	if err != nil {
		return "", errors.Wrapf(err, "cannot generate criteria '%s'", c.Name)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, c); err != nil {
		return "", errors.Wrapf(err, "cannot generate criteria '%s'", c.Name)
	}

	return buf.String(), nil
}

type CriteriaFieldType string

const (
	CriteriaFieldInt64  CriteriaFieldType = "int64"
	CriteriaFieldString CriteriaFieldType = "string"
)

func ParseCriteriaFieldType(t string) (CriteriaFieldType, error) {
	switch CriteriaFieldType(t) {
	case CriteriaFieldInt64:
		return CriteriaFieldInt64, nil
	case CriteriaFieldString:
		return CriteriaFieldString, nil
	}
	return "", errors.Errorf("invalid given criteria field type value '%s'", t)
}

type CriteriaField struct {
	Name string
	Type CriteriaFieldType
}
