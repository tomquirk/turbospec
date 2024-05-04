package turbospec

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
)

var TYPE_ALIAS_KEYWORD string = "type"

type TSType struct {
	Name  string
	Type  string
	Alias *string
}

type TSFormatter struct {
	TypeAliasTemplate   *template.Template
	DefaultTypeTemplate *template.Template
}

func NewTSFormatter() (*TSFormatter, error) {
	typeAliasTemplate, err := template.New("typeAliasTemplate").Parse("type {{.Name}} = {{.Type}};")
	if err != nil {
		return nil, err
	}
	defaultTypeTemplate, err := template.New("defaultTypeTemplate").Parse("{{.Name}}: {{.Type}};")
	if err != nil {
		return nil, err
	}

	return &TSFormatter{
		TypeAliasTemplate:   typeAliasTemplate,
		DefaultTypeTemplate: defaultTypeTemplate,
	}, err
}

func (f *TSFormatter) Format(tsType *TSType) (string, error) {
	typeTemplate := f.DefaultTypeTemplate
	if tsType.Alias != nil && *tsType.Alias == TYPE_ALIAS_KEYWORD {
		typeTemplate = f.TypeAliasTemplate
	}

	builder := strings.Builder{}
	err := typeTemplate.Execute(&builder, tsType)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (f *TSFormatter) ToTSPropertyObject(schema *openapi3.Schema) string {
	var properties []string
	for k, v := range schema.Properties {
		v.Value.Title = k
		property, err := f.ToTSType(v, nil)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}
		properties = append(properties, property)
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(properties, "\n"))
}

func (f *TSFormatter) ToTSPropertyArray(schema *openapi3.Schema) string {
	return "unknown[] // TODO fix"
}

func refToTypeName(ref string) string {
	tok := strings.Split(ref, "/")
	return tok[(len(tok) - 1)]
}

func (f *TSFormatter) ToTSType(schema *openapi3.SchemaRef, alias *string) (string, error) {
	tsType := TSType{
		Name:  schema.Value.Title, // TODO convert kebab-case to camelCase, add _ to names starting with a number
		Type:  "unknown // TODO fix",
		Alias: alias,
	}
	if schema.Ref != "" {
		tsType.Type = refToTypeName(schema.Ref)
	} else if schema.Value.Type.Is("string") {
		tsType.Type = "string"
	} else if schema.Value.Type.Is("integer") {
		tsType.Type = "number"
	} else if schema.Value.Type.Is("boolean") {
		tsType.Type = "boolean"
	} else if schema.Value.Type.Is("object") {
		tsType.Type = f.ToTSPropertyObject(schema.Value)
	} else if schema.Value.Type.Is("array") {
		tsType.Type = f.ToTSPropertyArray(schema.Value)
	}

	tsTypeStr, err := f.Format(&tsType)
	if err != nil {
		return "", err
	}

	return tsTypeStr, nil
}

func WriteTypes(spec *openapi3.T, writer io.Writer) {
	formatter, err := NewTSFormatter()
	if err != nil {
		log.Fatalln("Could not start formatter", err)
	}

	for k, v := range spec.Components.Schemas {
		v.Value.Title = k
		// set "root"-level types to use type alias (type x = { ... })
		typeStr, err := formatter.ToTSType(v, &TYPE_ALIAS_KEYWORD)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}
		writer.Write([]byte(typeStr))
	}
}

func LoadSpec(specFilePath string) (*openapi3.T, error) {
	fmt.Println(specFilePath)
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specFilePath)
	return doc, err
}
