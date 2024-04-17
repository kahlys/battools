package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"sort"
	"strings"
)

func main() {
	packagePath := flag.String("package", "", "package path")
	structName := flag.String("struct", "", "struct name")
	flag.Parse()

	check(*packagePath == "", "Package path is required")
	check(*structName == "", "Struct name is required")

	pkgs, err := parser.ParseDir(
		token.NewFileSet(),
		*packagePath,
		func(info os.FileInfo) bool {
			return !strings.HasSuffix(info.Name(), "_test.go")
		},
		parser.ParseComments,
	)
	checkErr(err, "Failed to parse package")
	check(len(pkgs) == 0, "No packages found")
	check(len(pkgs) > 1, "Multiple packages found")

	methods := []string{}
	for _, file := range firstPackage(pkgs).Files {
		methods = append(methods, getMethods(*structName, file)...)
	}

	res, err := generate(*structName, methods)
	if err != nil {
		checkErr(err, "Failed to generate interface")
	}

	fmt.Println(res)
}

func getMethods(structName string, file *ast.File) []string {
	var methods []string
	ast.Inspect(file, func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if funcDecl.Recv == nil {
			return true
		}
		for _, field := range funcDecl.Recv.List {
			starExpr, ok := field.Type.(*ast.StarExpr)
			if !ok {
				continue
			}
			ident, ok := starExpr.X.(*ast.Ident)
			if !ok {
				continue
			}
			if ident.Name == structName {
				params := make([]string, len(funcDecl.Type.Params.List))
				for i, param := range funcDecl.Type.Params.List {
					params[i] = fmt.Sprintf("%s %s", param.Names[0].Name, typeExprToString(param.Type))

				}

				results := []string{}
				if funcDecl.Type.Results != nil {
					results = make([]string, len(funcDecl.Type.Results.List))
					for i, result := range funcDecl.Type.Results.List {
						results[i] = typeExprToString(result.Type)

					}
				}

				method := fmt.Sprintf("%s(%s) (%s)", funcDecl.Name.Name, strings.Join(params, ", "), strings.Join(results, ", "))
				methods = append(methods, method)
			}
		}
		return true
	})
	return methods
}

func generate(structName string, methods []string) (string, error) {
	buf := bytes.Buffer{}

	sort.Strings(methods)

	fmt.Fprintf(&buf, "type %s interface {\n", structName)
	for _, method := range methods {
		fmt.Fprintf(&buf, "\t%s\n", method)
	}
	fmt.Fprintf(&buf, "}")

	pretty, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(pretty), nil
}

func typeExprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), expr)
	return buf.String()
}

func firstPackage(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

func checkErr(err error, msg string) {
	if err != nil {
		fmt.Printf("%s\n%s\n", msg, err)
		os.Exit(1)
	}
}

func check(b bool, msg string) {
	if b {
		fmt.Println(msg)
		os.Exit(1)
	}
}
