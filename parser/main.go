package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	typeparser "github.com/fsamin/go-typeparser"
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
	ts, err := typeparser.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "cannot parse file at '%s'", path)
	}

	var crs []*ultimate.Criteria

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

	for _, cr := range crs {
		src, err := cr.Generate()
		if err != nil {
			return err
		}
		fmt.Println(src)
	}

	return nil
}

func extractCriteria(criteriaType string, t typeparser.Type) (*ultimate.Criteria, error) {
	ct, err := ultimate.ParseCriteriaType(criteriaType)
	if err != nil {
		return nil, err
	}

	c := ultimate.NewCriteria(t.Name(), ct)

	for _, field := range t.Fields() {
		opts := field.TagValue("ultimate")
		if len(opts) < 2 || opts[1] != "criteria" {
			continue
		}

		fieldType, err := ultimate.ParseCriteriaFieldType(field.Type())
		if err != nil {
			return nil, err
		}
		c.AddField(opts[0], fieldType)
	}

	return &c, nil
}
