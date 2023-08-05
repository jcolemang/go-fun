package main

import (
	"os"
	"fmt"

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
		prog, err := Compile(ast)
		ctx.FatalIfErrorf(err)
		fmt.Println("Final program")
		fmt.Println(ArmProgramToString(prog))

		err = CompileToFile(ast, "assembly.s")
		ctx.FatalIfErrorf(err)
	}
}
