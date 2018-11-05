package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type criteriaType int

const (
	criteriaSQL criteriaType = iota
	criteriaMongo
)

func parseCriteriaType(t string) (criteriaType, error) {
	switch t {
	case "sql":
		return criteriaSQL, nil
	case "mongo":
		return criteriaMongo, nil
	}
	return 0, errors.Errorf("invalid given criteria type value '%s'", t)
}

func newCriteria(name string, t criteriaType) criteria {
	return criteria{Name: name, Type: t}
}

type criteria struct {
	Package, Name string
	Type          criteriaType
	Fields        []criteriaField
}

func (c *criteria) addField(name string, fieldType criteriaFieldType) {
	c.Fields = append(c.Fields, criteriaField{Name: name, Type: fieldType})
}

func (c criteria) generate() (string, error) {
	content := `
func NewCriteria{{.Name | Title}}() *Criteria{{.Name | Title}} {
	return &Criteria{{.Name | Title}}{}
}

type Criteria{{.Name | Title}} struct { 
	{{- range .Fields}}
	{{.Name}}s []{{.Type}}
	{{- end}}
}
{{range .Fields}}
func (c *Criteria{{$.Name | Title}}) {{.Name | Title }}s(v ...{{.Type}}) *Criteria{{$.Name | Title}}{
	c.{{.Name}}s = v
	return c
}
{{end -}}
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

type criteriaFieldType string

const (
	criteriaFieldInt64  criteriaFieldType = "int64"
	criteriaFieldString criteriaFieldType = "string"
)

func parseCriteriaFieldType(t string) (criteriaFieldType, error) {
	switch criteriaFieldType(t) {
	case criteriaFieldInt64:
		return criteriaFieldInt64, nil
	case criteriaFieldString:
		return criteriaFieldString, nil
	}
	return "", errors.Errorf("invalid given criteria field type value '%s'", t)
}

type criteriaField struct {
	Name string
	Type criteriaFieldType
}
