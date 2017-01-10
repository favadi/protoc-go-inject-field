package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var testInputFile = "pb/test.pb.go"

func TestFieldFromComment(t *testing.T) {
	var tests = []struct {
		comment string
		result  *customField
	}{
		{
			comment: `//@inject_field: name string`,
			result:  &customField{fieldName: "name", fieldType: "string"}},
		{
			comment: `//   @inject_field: age int`,
			result:  &customField{fieldName: "age", fieldType: "int"}},
		{
			comment: `//@inject_field:    age    int32`,
			result:  &customField{fieldName: "age", fieldType: "int32"}},
		{
			comment: `//not_a_match `,
			result:  nil},
		{
			comment: `//@inject_field: `,
			result:  nil},
	}
	for _, test := range tests {
		result := fieldFromComment(test.comment)
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("wrong customField extraction from customField: %q", test.comment)
		}
	}
}

func TestInjectField(t *testing.T) {
	injected := injectField(
		[]byte(`type abcdef struct {
				xyz string
}
`),
		textArea{
			start:     1,
			end:       38,
			insertPos: 36,
			fields: []*customField{
				{fieldName: "a", fieldType: "b"},
				{fieldName: "c", fieldType: "d"},
			},
		})
	expected :=
		[]byte(`type abcdef struct {
				xyz string

	// custom fields
	a b
	c d
}

// custom fields getter/setter
func (m *) A() b {
	return m.a
}
func (m *) SetA(in b){
	m.a = in
}
func (m *) C() d {
	return m.c
}
func (m *) SetC(in d){
	m.c = in
}
`)
	if string(injected) != string(expected) {
		t.Fatal("custom fields are not injected properly")
	}
}

func TestParseWriteFile(t *testing.T) {
	// copy the pb.go to a temp file
	tmpFile, err := ioutil.TempFile("", "protoc-go-inject-field-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	input, err := os.Open(testInputFile)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()

	if _, err = io.Copy(tmpFile, input); err != nil {
		t.Fatal(err)
	}

	if err = tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	areas, err := parseFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(areas) != 2 {
		t.Fatalf("expected %d structs to inject custom fields", 2)
	}

	if len(areas[0].fields) != 2 {
		t.Fatal("expected 2 custom fields for Human struct")
	}
	if len(areas[1].fields) != 1 {
		t.Fatal("expected 1 custom fields for Robot struct")
	}

	if err = writeFile(tmpFile.Name(), areas); err != nil {
		t.Fatal(err)
	}

	injected, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expectedSnippets := [][]byte{
		[]byte(`// custom fields
	age int
	spouse *Human
`),
		[]byte(`// custom fields getter/setter
func (m *Human) Age() int {
	return m.age
}
func (m *Human) SetAge(in int){
	m.age = in
}
func (m *Human) Spouse() *Human {
	return m.spouse
}
func (m *Human) SetSpouse(in *Human){
	m.spouse = in
}`),
		[]byte(`// custom fields
	model string`),
		[]byte(`// custom fields getter/setter
func (m *Robot) Model() string {
	return m.model
}
func (m *Robot) SetModel(in string){
	m.model = in
}`),
	}
	for _, snippet := range expectedSnippets {
		if !bytes.Contains(injected, snippet) {
			t.Errorf("injected file does not contain: %q", snippet)
		}
	}
}
