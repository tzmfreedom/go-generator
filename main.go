package main

import (
	"io"
	"os"

	"github.com/tzmfreedom/go-generator/generator"
)

const defaultVersion = "7.4"

func main() {
	err := Main()
	if err != nil {
		panic(err)
	}
}

func Main() error {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	version := defaultVersion
	if os.Getenv("PHP_VERSION") != "" {
		version = os.Getenv("PHP_VERSION")
	}
	return generator.Generate(buf, version, os.Getenv("DEBUG") != "")
}
