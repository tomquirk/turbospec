// Package builder implements [turbospec] `Builder`s that generate
// strings from Typescript types.
//
// Builders are mostly intended to generate code
package builder

import (
	"strings"
	"text/template"

	"github.com/tomquirk/turbospec/pkg/ts"
)

var TYPE_ALIAS_KEYWORD string = "type"

// A TSBuilder generates Typescript types.
type TSBuilder struct {
	defaultTypeTemplate *template.Template
	typeAliasTemplate   *template.Template
}

// NewTSBuilder returns a new TSBuilder.
func NewTSBuilder() (*TSBuilder, error) {
	typeAliasTemplate, err := template.New("typeAliasTemplate").Parse("type {{.Name}} = {{.Type}};")
	if err != nil {
		return nil, err
	}
	defaultTypeTemplate, err := template.New("defaultTypeTemplate").Parse("{{.Name}}{{if .Required}}{{else}}?{{end}}: {{.Type}};")
	if err != nil {
		return nil, err
	}

	return &TSBuilder{
		defaultTypeTemplate: defaultTypeTemplate,
		typeAliasTemplate:   typeAliasTemplate,
	}, nil
}

// Build returns Typescript code (as a string) for given tsType.
//
// If tsType.Type is empty, the `unknown` type is used.
func (f TSBuilder) Build(tsType ts.TSType) (string, error) {
	if tsType.Type == "" {
		tsType.Type = "unknown"
	}

	typeTemplate := f.defaultTypeTemplate
	if tsType.Alias == TYPE_ALIAS_KEYWORD {
		typeTemplate = f.typeAliasTemplate
	}

	builder := strings.Builder{}
	err := typeTemplate.Execute(&builder, tsType)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}
