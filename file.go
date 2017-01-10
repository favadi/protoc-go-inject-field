package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// parseFile scans the input file and returns list of structs to inject custom fields to.
func parseFile(inputPath string) (areas []textArea, err error) {
	areas = []textArea{}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inputPath, nil, parser.ParseComments)
	if err != nil {
		return
	}

	for _, decl := range f.Decls {
		// check if is generic declaration
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		var typeSpec *ast.TypeSpec
		for _, spec := range genDecl.Specs {
			if ts, tsOK := spec.(*ast.TypeSpec); tsOK {
				typeSpec = ts
				break
			}
		}

		// skip if can't get type spec
		if typeSpec == nil {
			continue
		}

		// not a struct, skip
		structDecl, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		area := textArea{
			name:      typeSpec.Name.String(),
			start:     int(typeSpec.Pos()),
			end:       int(typeSpec.End()),
			insertPos: int(structDecl.Fields.Closing) - 1,
			fields:    []*customField{},
		}

		if genDecl.Doc == nil {
			continue
		}

		// build the list of text areas from comments
		for _, comment := range genDecl.Doc.List {
			field := fieldFromComment(comment.Text)

			if field == nil || len(field.fieldName) == 0 || len(field.fieldName) == 0 {
				continue
			}

			// only inject private fields
			firtChar := string(field.fieldName[0])
			if strings.ToUpper(firtChar) == firtChar {
				continue
			}

			area.fields = append(area.fields, field)
		}
		areas = append(areas, area)
	}
	return
}

// writeFile updates the given files with given text areas.
func writeFile(inputPath string, areas []textArea) (err error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return
	}

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	if err = f.Close(); err != nil {
		return
	}

	// inject custom fields from the end of file first to preserve order
	for i := len(areas) - 1; i >= 0; i-- {
		contents = injectField(contents, areas[i])
	}
	if err = ioutil.WriteFile(inputPath, contents, 0644); err != nil {
		return
	}

	log.Printf("file %q is injected with custom fields", inputPath)
	return
}
