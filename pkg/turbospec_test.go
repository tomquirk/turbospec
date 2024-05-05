package turbospec

import (
	"testing"
)

func TestLoadSpec(t *testing.T) {
	specFilePath := "../data/test_openapi_spec.yaml"
	want := "3.0.0"

	transformer, err := NewOpenapiTsTransformer()
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = transformer.LoadSpec(specFilePath)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if transformer.openapiDoc == nil || transformer.openapiDoc.OpenAPI != want {
		t.Fatalf("doc's OpenAPI prop didn't match expected: %s", want)
	}
}
