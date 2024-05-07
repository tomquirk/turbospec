package turbospec

import (
	"fmt"
	"io"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
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

func (t *OpenapiTsTransformer) ToTSPropertyObject(schema *openapi3.Schema) string {
	var properties []string
	for k, v := range schema.Properties {
		v.Value.Title = k
		property, err := t.ToTSType(v, "")
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}
		properties = append(properties, property)
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(properties, "\n"))
}

func (t *OpenapiTsTransformer) ToTSPropertyArray(schema *openapi3.Schema) string {
	itemType := "unknown"
	if schema.Items.Value.Type.Is("string") {
		itemType = "string"
	} else if schema.Items.Value.Type.Is("integer") {
		itemType = "number"
	} else if schema.Items.Value.Type.Is("boolean") {
		itemType = "boolean"
	} else if schema.Items.Value.Type.Is("object") {
		itemType = t.ToTSPropertyObject(schema.Items.Value)
	} else if schema.Items.Value.Type.Is("array") {
		itemType = t.ToTSPropertyArray(schema.Items.Value)
	}

	return fmt.Sprintf("%s[]", itemType)
}

func (t *OpenapiTsTransformer) ToTSType(schema *openapi3.SchemaRef, alias string) (string, error) {
	tsType := ts.TSType{
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
		tsType.Type = t.ToTSPropertyObject(schema.Value)
	} else if schema.Value.Type.Is("array") {
		tsType.Type = t.ToTSPropertyArray(schema.Value)
	}

	tsTypeStr, err := t.tsBuilder.Build(tsType)
	if err != nil {
		return "", err
	}

	return tsTypeStr, nil
}

func (t *OpenapiTsTransformer) Transform(openapiDoc *openapi3.T, out io.Writer) {
	for k, v := range openapiDoc.Components.Schemas {
		v.Value.Title = k
		// set "root"-level types to use type alias (type x = { ... })
		typeStr, err := t.ToTSType(v, ts.TYPE_ALIAS_KEYWORD)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}

		out.Write([]byte(typeStr))
	}
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

func refToTypeName(ref string) string {
	tok := strings.Split(ref, "/")
	return tok[(len(tok) - 1)]
}
