package main

import (
    "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Main Language
type Program struct {
	Expr *Expr `@@`
}

type Expr struct {
	Num *Num     `@@`
    Var *Var     `| @@`
	Let *LetExpr `| @@`
    App []*Expr  `| "(" @@ @@* ")" `
}

type LetExpr struct {
	LetAssignments []*Assignment `"(" "let" "(" @@ @@* ")"`
	LetBody *Expr                   `@@ ")"`
}

type Assignment struct {
	Ref *Var   `"(" @@`
	Expr *Expr `@@ ")"`
}

type Num struct {
	Value *float64 `@Int`
}

type Var struct {
    Name *string `@Ident`
    Temp *int
}

func GetLanguageParser() *participle.Parser[Program] {
	basicLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `(?i)rem[^\n]*`},
		{"String", `"(\\"|[^"])*"`},
		{"Ident", `[a-zA-Z_\-+]\w*`},
		{"Int", `[-]?(\d+)`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{"whitespace", `[ \t\n\r]+`},
	})
    parser := participle.MustBuild[Program](
        participle.Lexer(basicLexer),
    )
    return parser
}
