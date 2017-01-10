package main

import (
	"fmt"
	"strings"
)

// injectField injects the custom field to the struct.
// TODO: we should modify ast tree instead of manipulating bytes directly, but it will need more research
// TODO: we should use text/template instead of fmt.Sprintf
func injectField(contents []byte, area textArea) (injected []byte) {
	customFieldsContents := []byte("\n\t// custom fields\n")
	for _, field := range area.fields {
		customFieldsContents = append(
			customFieldsContents,
			[]byte(fmt.Sprintf("\t%s %s\n", field.fieldName, field.fieldType))...)
	}

	helperMethodContents := []byte("\n// custom fields getter/setter\n")
	for _, field := range area.fields {
		helperMethodContents = append(
			helperMethodContents,
			[]byte(fmt.Sprintf(
				`func (m *%s) %s() %s {
	return m.%s
}
func (m *%s) Set%s(in %s){
	m.%s = in
}
`,
				// getter params
				area.name,
				strings.Title(field.fieldName),
				field.fieldType,
				field.fieldName,
				// setter params
				area.name,
				strings.Title(field.fieldName),
				field.fieldType,
				field.fieldName))...)
	}

	injected = append(
		contents[:area.end],
		append(helperMethodContents, contents[area.end:]...)...)

	injected = append(
		injected[:area.insertPos],
		append(customFieldsContents, injected[area.insertPos:]...)...)
	return
}

// fieldFromComment gets the custom customField information from comment.
func fieldFromComment(comment string) *customField {
	match := rComment.FindStringSubmatch(comment)
	if len(match) == 3 {
		return &customField{fieldName: match[1], fieldType: match[2]}
	}
	return nil
}
