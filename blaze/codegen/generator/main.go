package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"
)

func generate(value string) (string, error) {
	var buf bytes.Buffer

	var templ = template.Must(template.New("gen").Parse(`
		package codegen

		import "fmt"

		func Hello() string {
			return fmt.Sprintf("Hello, %s!", "{{.}}")
		}`,
	))

	if err := templ.Execute(&buf, value); err != nil {
		return "", err
	}

	pretty, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(pretty), nil
}

func main() {
	var value = flag.String("name", "world", "name used in the generated hello package")
	var out = flag.String("out", "", "output file name")

	flag.Parse()

	s, err := generate(*value)
	if err != nil {
		check(err)
	}

	var wr io.Writer = os.Stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			check(err)
		}
		defer f.Close()
		wr = f
	}

	_, err = fmt.Fprint(wr, s)
	if err != nil {
		check(err)
	}
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
