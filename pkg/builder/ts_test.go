package builder

import (
	"testing"

	"github.com/tomquirk/turbospec/pkg/ts"
)

func TestNewTSBuilder(t *testing.T) {
	tsBuilder, err := NewTSBuilder()
	if err != nil || tsBuilder.defaultTypeTemplate == nil || tsBuilder.typeAliasTemplate == nil {
		t.Fatalf("Couldn't create TSBuilder")
	}
}

func TestBuild(t *testing.T) {
	tsBuilder, err := NewTSBuilder()
	if err != nil || tsBuilder.defaultTypeTemplate == nil || tsBuilder.typeAliasTemplate == nil {
		t.Fatalf("Couldn't create TSBuilder")
	}

	type test struct {
		tsType ts.TSType
		want   string
	}

	tests := []test{
		{tsType: ts.TSType{Name: "user"}, want: "user: unknown;"},
		{tsType: ts.TSType{Name: "user", Type: "string"}, want: "user: string;"},
		{tsType: ts.TSType{Name: "user", Type: "number"}, want: "user: number;"},
		{tsType: ts.TSType{Name: "user", Type: "number[]"}, want: "user: number[];"},
		{tsType: ts.TSType{Name: "User", Type: "string", Alias: TYPE_ALIAS_KEYWORD}, want: "type User = string;"},
		{tsType: ts.TSType{Name: "user", Type: "string", Alias: TYPE_ALIAS_KEYWORD}, want: "type user = string;"},
	}

	for _, tc := range tests {
		out, err := tsBuilder.Build(tc.tsType)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if out != tc.want {
			t.Fatalf("%s != %s", out, tc.want)
		}
	}

}
