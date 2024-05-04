package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	turbospec "github.com/tomquirk/turbospec/pkg"
)

func run() error {
	args := os.Args
	specFilePath := args[1]
	if specFilePath == "" {
		return errors.New("no spec file path specified")
	}
	doc, err := turbospec.LoadSpec(specFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("Doc loaded, openapi=%s\n", doc.OpenAPI)

	builder := strings.Builder{}
	turbospec.WriteTypes(doc, &builder)

	fmt.Println(builder.String())

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
