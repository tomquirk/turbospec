package turbospec

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/tomquirk/turbospec/pkg/builder"
	"github.com/tomquirk/turbospec/pkg/ts"
)

type TSBuilder interface {
	Build(tsType ts.TSType) (string, error)
}

type OpenapiTsTransformer struct {
	tsBuilder TSBuilder
}

type tsTypeOptions struct {
	optional bool
	alias    string
}

type tsFormatOptions struct {
	tabs int8
}

func NewOpenapiTsTransformer() (*OpenapiTsTransformer, error) {
	tsBuilder, err := builder.NewTSBuilder()
	if err != nil {
		return nil, err
	}

	return &OpenapiTsTransformer{
		tsBuilder: tsBuilder,
	}, nil
}

func LoadSpec(specFilePath string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specFilePath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Doc loaded, openapi=%s\n", doc.OpenAPI)

	return doc, nil
}

func (t *OpenapiTsTransformer) Transform(openapiDoc *openapi3.T, out io.Writer) {
	for k, v := range openapiDoc.Components.Schemas {
		v.Value.Title = k
		// set "root"-level types to use type alias (type MyType = { ... })
		typeStr, err := t.ToTSType(v, tsTypeOptions{
			optional: false, alias: ts.TYPE_ALIAS_KEYWORD,
		}, tsFormatOptions{1})
		if err != nil {
			log.Println(err.Error())
			continue
		}

		out.Write([]byte("\n" + typeStr))
	}
}

func (t *OpenapiTsTransformer) ToTSType(schema *openapi3.SchemaRef, typeOpts tsTypeOptions, formatOpts tsFormatOptions) (string, error) {
	tsType := ts.TSType{
		Name:     normalizeTypeName(schema.Value.Title),
		Type:     "unknown // TODO fix",
		Alias:    typeOpts.alias,
		Optional: typeOpts.optional,
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
		tsType.Type = t.ToTSPropertyObject(schema.Value, formatOpts)
	} else if schema.Value.Type.Is("array") {
		tsType.Type = t.ToTSPropertyArray(schema.Value, formatOpts)
	}

	tsTypeStr, err := t.tsBuilder.Build(tsType)
	if err != nil {
		return "", err
	}

	return tsTypeStr, nil
}

func (t *OpenapiTsTransformer) ToTSPropertyObject(schema *openapi3.Schema, formatOpts tsFormatOptions) string {
	tabs := strings.Repeat("\t", int(formatOpts.tabs))
	tabsClosingBrace := strings.Repeat("\t", int(formatOpts.tabs-1))

	var properties []string
	for k, v := range schema.Properties {
		v.Value.Title = k
		required := slices.Contains(schema.Required, k)
		fmt.Println("optional", !required)
		property, err := t.ToTSType(v, tsTypeOptions{optional: !required, alias: ""}, tsFormatOptions{tabs: formatOpts.tabs + 1})
		if err != nil {
			log.Println(err.Error())
			continue
		}
		property = tabs + property
		properties = append(properties, property)
	}

	return fmt.Sprintf("{\n%s\n%s}", strings.Join(properties, "\n"), tabsClosingBrace)
}

func (t *OpenapiTsTransformer) ToTSPropertyArray(schema *openapi3.Schema, formatOpts tsFormatOptions) string {
	itemType := "unknown"
	if schema.Items.Value.Type.Is("string") {
		itemType = "string"
	} else if schema.Items.Value.Type.Is("integer") {
		itemType = "number"
	} else if schema.Items.Value.Type.Is("boolean") {
		itemType = "boolean"
	} else if schema.Items.Value.Type.Is("object") {
		itemType = t.ToTSPropertyObject(schema.Items.Value, formatOpts)
	} else if schema.Items.Value.Type.Is("array") {
		itemType = t.ToTSPropertyArray(schema.Items.Value, formatOpts)
	}

	return fmt.Sprintf("%s[]", itemType)
}

func normalizeTypeName(name string) string {
	if name[0] >= '0' && name[0] <= '9' {
		name = "_" + name
	}

	if !strings.Contains(name, "-") {
		return name
	}
	return strcase.ToLowerCamel(name)
}

func refToTypeName(ref string) string {
	tok := strings.Split(ref, "/")
	return tok[(len(tok) - 1)]
}
