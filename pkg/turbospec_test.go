package turbospec

import (
	"testing"
)

func TestLoadSpec(t *testing.T) {
	specFilePath := "../data/test_openapi_spec.yaml"
	want := "3.0.0"

	doc, err := LoadSpec(specFilePath)

	if doc == nil {
		t.Fatalf("No doc loaded")
	}
	if err != nil {
		t.Fatalf(err.Error())
	}
	if doc.OpenAPI != want {
		t.Fatalf("doc's OpenAPI prop didn't match expected: %s", want)
	}
}

func BenchmarkLoadSpec(b *testing.B) {
	specFilePath := "../data/test_openapi_spec.json"
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		LoadSpec(specFilePath)
	}
}
