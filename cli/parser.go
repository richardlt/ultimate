package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	typeparser "github.com/fsamin/go-typeparser"
	"github.com/pkg/errors"
)

func parse(path, dest, packageName string) error {
	ts, err := typeparser.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "cannot parse file at '%s'", path)
	}

	var crs []*criteria

	for _, t := range ts {
		var commentCriteria string
		for _, d := range t.Docs() {
			if strings.Contains(d, "ultimate") {
				commentCriteria = d
				break
			}
		}
		if commentCriteria != "" {
			split := strings.Split(commentCriteria, ":")
			if len(split) == 2 {
				cr, err := extractCriteria(split[1], t)
				if err != nil {
					return errors.Wrapf(err, "cannot generate criteria for struct '%s'", t.Name())
				}
				crs = append(crs, cr)
			}
		}
	}

	f := newFile(packageName)
	for _, cr := range crs {
		src, err := cr.generate()
		if err != nil {
			return err
		}
		f.Criteria = append(f.Criteria, src)
	}
	src, err := f.generate()
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(path)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	filePath := fmt.Sprintf("%s/%s.criteria.go", dest, fileName)
	if err := ioutil.WriteFile(filePath, []byte(src), os.ModePerm); err != nil {
		return errors.Wrapf(err, "cannot write file for criteria at %s", dest)
	}

	return nil
}

func extractCriteria(criteriaType string, t typeparser.Type) (*criteria, error) {
	ct, err := parseCriteriaType(criteriaType)
	if err != nil {
		return nil, err
	}

	c := newCriteria(t.Name(), ct)

	for _, field := range t.Fields() {
		opts := field.TagValue("ultimate")
		if len(opts) < 2 || opts[1] != "criteria" {
			continue
		}

		fieldType, err := parseCriteriaFieldType(field.Type())
		if err != nil {
			return nil, err
		}
		c.addField(opts[0], fieldType)
	}

	return &c, nil
}
