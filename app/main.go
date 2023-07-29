package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/repr"
)

// Code to actually run
var cli struct {
	Files []string `arg:"" optional:"" type:"existingfile" help:"GraphQL schema files to parse."`
}

func main() {
    parser := GetLanguageParser()
	ctx := kong.Parse(&cli)
	for _, file := range cli.Files {
		r, err := os.Open(file)
		ctx.FatalIfErrorf(err)
		ast, err := parser.Parse(file, r)
		r.Close()
		repr.Println(ast)
		ctx.FatalIfErrorf(err)

		newProg, err := Compile(ast)
		repr.Println(newProg)
		ctx.FatalIfErrorf(err)
	}
}
