package main

import (
	"errors"
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

	f, err := os.Create("types.ts")
	if err != nil {
		return err
	}
	// builder := strings.Builder{}
	// fmt.Println(builder.String())

	transformer, err := turbospec.NewOpenapiTsTransformer()
	if err != nil {
		return err
	}

	spec, err := turbospec.LoadSpec(specFilePath)
	if err != nil {
		return err
	}

	transformer.Transform(spec, f)

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
