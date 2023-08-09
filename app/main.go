package main

import (
	"os"
    "fmt"
    "language/pkg/languages"

	"github.com/alecthomas/kong"
)

// Code to actually run
var cli struct {
	InputFile string `arg:"" optional:"" type:"existingfile" help:"GraphQL schema files to parse."`
	OutputFile string `arg:"" optional:"" type:"newfile" help:"GraphQL schema files to parse."`
}

func main() {
    parser := languages.GetLanguageParser()
	ctx := kong.Parse(&cli)
    file, result := cli.InputFile, cli.OutputFile
    r, err := os.Open(file)
    ctx.FatalIfErrorf(err)
    ast, err := parser.Parse(file, r)
    r.Close()
    ctx.FatalIfErrorf(err)
    fmt.Println(languages.ProgToString(ast))

    debug := true
    err = CompileToFile(ast, result, debug)
    ctx.FatalIfErrorf(err)
}
