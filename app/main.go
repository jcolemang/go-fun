package main

import (
	"os"

	"github.com/alecthomas/kong"
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
		ctx.FatalIfErrorf(err)

        debug := false
		err = CompileToFile(ast, "assembly.s", debug)
		ctx.FatalIfErrorf(err)
	}
}
