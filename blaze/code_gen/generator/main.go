package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"text/template"
)

func generate(value string) (string, error) {
	var buf bytes.Buffer

	var templ = template.Must(template.New("gen").Parse(`
		package hello

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

	flag.Parse()

	s, err := generate(*value)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(s)
}
