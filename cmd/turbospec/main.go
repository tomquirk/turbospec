package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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

	// builder := strings.Builder{}
	// turbospec.WriteTypes(doc, &builder)
	// fmt.Println(builder.String())

	f, err := os.Create("types.ts")
	if err != nil {
		return err
	}
	turbospec.WriteTypes(doc, f)

	return nil
}

func main() {
	start := time.Now()
	err := run()
	elapsed := time.Since(start)
	log.Printf("%s", elapsed)
	if err != nil {
		log.Fatal(err.Error())
	}
}
