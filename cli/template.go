package main

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"Title":   strings.Title,
	"ToLower": strings.ToLower,
}
