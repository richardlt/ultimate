package main

import (
	"log"
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "ultimate"
	app.Commands = []cli.Command{
		{
			Name:      "criteria",
			ArgsUsage: "first.go second.go third.go",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "output-dir", Value: "."},
				cli.StringFlag{Name: "go-package", Value: "criteria"},
			},
			Action: func(c *cli.Context) error {
				for _, a := range c.Args() {
					if err := parse(a, c.String("output-dir"),
						c.String("go-package")); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
