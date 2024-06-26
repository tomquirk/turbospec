package turbospec

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestLoadSpec(t *testing.T) {
	specFilePath := "../data/test_openapi_spec.yaml"
	want := "3.0.0"

	spec, err := LoadSpec(specFilePath)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if spec == nil || spec.OpenAPI != want {
		t.Fatalf("doc's OpenAPI prop didn't match expected: %s", want)
	}
}

func TestToTSType(t *testing.T) {

	transformer, err := NewOpenapiTsTransformer()
	if err != nil {
		t.Fatalf(err.Error())
	}

	type test struct {
		schemaRef *openapi3.SchemaRef
		alias     string
		want      string
	}
	testCases := []test{
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"integer"},
					Title: "age",
				},
			},
			alias: "",
			want:  "age: number;",
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"boolean"},
					Title: "isAlive",
				},
			},
			alias: "",
			want:  "isAlive: boolean;",
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"string"},
					Title: "name",
				},
			},
			alias: "",
			want:  "name: string;",
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"string"},
					Title: "Name",
				},
			},
			alias: "type",
			want:  "type Name = string;",
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"array"},
					Title: "attributes",
					Items: &openapi3.SchemaRef{
						Ref: "",
						Value: &openapi3.Schema{
							Type: &openapi3.Types{"string"},
						},
					},
				},
			},
			alias: "",
			want:  "attributes: string[];",
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"array"},
					Title: "attributes",
					Items: &openapi3.SchemaRef{
						Ref: "",
						Value: &openapi3.Schema{
							Type: &openapi3.Types{"integer"},
						},
					},
				},
			},
			alias: "",
			want:  "attributes: number[];",
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"array"},
					Title: "attributes",
					Items: &openapi3.SchemaRef{
						Ref: "",
						Value: &openapi3.Schema{
							Type: &openapi3.Types{"boolean"},
						},
					},
				},
			},
			alias: "",
			want:  "attributes: boolean[];",
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"array"},
					Title: "attributes",
					Items: &openapi3.SchemaRef{
						Ref: "",
						Value: &openapi3.Schema{
							Type:     &openapi3.Types{"object"},
							Required: []string{"name"},
							Properties: openapi3.Schemas{
								"name": &openapi3.SchemaRef{
									Ref: "",
									Value: &openapi3.Schema{
										Type:  &openapi3.Types{"string"},
										Title: "name",
									},
								},
							},
						},
					},
				},
			},
			alias: "",
			want: `attributes: {
	name: string;
}[];`,
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"array"},
					Title: "attributes",
					Items: &openapi3.SchemaRef{
						Ref: "",
						Value: &openapi3.Schema{
							Type: &openapi3.Types{"array"},
							Items: &openapi3.SchemaRef{
								Ref: "",
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"string"},
								},
							},
						},
					},
				},
			},
			alias: "",
			want:  `attributes: string[][];`,
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:     &openapi3.Types{"object"},
					Title:    "Person",
					Required: []string{"name"},
					Properties: openapi3.Schemas{
						"name": &openapi3.SchemaRef{
							Ref: "",
							Value: &openapi3.Schema{
								Type:  &openapi3.Types{"string"},
								Title: "name",
							},
						},
					},
				},
			},
			alias: "type",
			want: `type Person = {
	name: string;
};`,
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"object"},
					Title: "Person",
					Properties: openapi3.Schemas{
						"name": &openapi3.SchemaRef{
							Ref: "",
							Value: &openapi3.Schema{
								Type:  &openapi3.Types{"string"},
								Title: "name",
							},
						},
					},
				},
			},
			alias: "type",
			want: `type Person = {
	name?: string;
};`,
		},
		{
			schemaRef: &openapi3.SchemaRef{
				Ref: "",
				Value: &openapi3.Schema{
					Type:     &openapi3.Types{"object"},
					Title:    "Person",
					Required: []string{"name"},
					Properties: openapi3.Schemas{
						"name": &openapi3.SchemaRef{
							Ref: "",
							Value: &openapi3.Schema{
								Type:  &openapi3.Types{"object"},
								Title: "name",
								Properties: openapi3.Schemas{
									"firstName": &openapi3.SchemaRef{
										Ref: "",
										Value: &openapi3.Schema{
											Type:  &openapi3.Types{"string"},
											Title: "firstName",
										},
									},
								},
							},
						},
					},
				},
			},
			alias: "type",
			want: `type Person = {
	name: {
		firstName?: string;
	};
};`,
		},
	}

	for _, tc := range testCases {
		res, err := transformer.ToTSType(tc.schemaRef, tsTypeOptions{false, tc.alias}, tsFormatOptions{1})
		if err != nil {
			t.Fatalf(err.Error())
		}

		if res != tc.want {
			t.Fatalf("got:\n%s, but want:\n%s", res, tc.want)
		}
	}
}
