package turbospec

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/tomquirk/turbospec/pkg/builder"
	"github.com/tomquirk/turbospec/pkg/ts"
)

type Builder interface {
	Build(tsType ts.TSType) (string, error)
}

type OpenapiTsTransformer struct {
	tsBuilder Builder
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
		typeStr, err := t.ToTSType(v, ts.TYPE_ALIAS_KEYWORD, 1, false)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}

		out.Write([]byte("\n" + typeStr))
	}
}

func (t *OpenapiTsTransformer) ToTSType(schema *openapi3.SchemaRef, alias string, distanceFromRoot int8, required bool) (string, error) {
	tsType := ts.TSType{
		Name:     normalizeTypeName(schema.Value.Title),
		Type:     "unknown // TODO fix",
		Alias:    alias,
		Required: required,
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
		tsType.Type = t.ToTSPropertyObject(schema.Value, distanceFromRoot)
	} else if schema.Value.Type.Is("array") {
		tsType.Type = t.ToTSPropertyArray(schema.Value, distanceFromRoot)
	}

	tsTypeStr, err := t.tsBuilder.Build(tsType)
	if err != nil {
		return "", err
	}

	return tsTypeStr, nil
}

func (t *OpenapiTsTransformer) ToTSPropertyObject(schema *openapi3.Schema, distanceFromRoot int8) string {
	tabs := strings.Repeat("\t", int(distanceFromRoot))
	tabsClosingBrace := strings.Repeat("\t", int(distanceFromRoot-1))

	var properties []string
	for k, v := range schema.Properties {
		v.Value.Title = k
		required := slices.Contains(schema.Required, k)
		property, err := t.ToTSType(v, "", distanceFromRoot+1, required)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}
		property = tabs + property
		properties = append(properties, property)
	}

	return fmt.Sprintf("{\n%s\n%s}", strings.Join(properties, "\n"), tabsClosingBrace)
}

func (t *OpenapiTsTransformer) ToTSPropertyArray(schema *openapi3.Schema, distanceFromRoot int8) string {
	itemType := "unknown"
	if schema.Items.Value.Type.Is("string") {
		itemType = "string"
	} else if schema.Items.Value.Type.Is("integer") {
		itemType = "number"
	} else if schema.Items.Value.Type.Is("boolean") {
		itemType = "boolean"
	} else if schema.Items.Value.Type.Is("object") {
		itemType = t.ToTSPropertyObject(schema.Items.Value, distanceFromRoot)
	} else if schema.Items.Value.Type.Is("array") {
		itemType = t.ToTSPropertyArray(schema.Items.Value, distanceFromRoot)
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
