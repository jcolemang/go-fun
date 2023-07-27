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
	NumVal *Num      `@@`
    VarExpr *Var     `| @@`
	LetExpr *LetExpr `| @@`
    AppExpr []*Expr  `| "(" @@ @@* ")" `
}

type LetExpr struct {
	LetAssignments []*LetAssignment `"(" "let" "(" @@ @@* ")"`
	LetBody *Expr                   `@@ ")"`
}

type LetAssignment struct {
	Ref *Var  `"(" @@`
	Val *Expr `@@ ")"`
}

type Num struct {
	Value *float64 `@Int`
}

type Var struct {
    Value *string `@Ident`
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