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
	Value int `@Int`
}

type Var struct {
    Name string `@Ident`
    Temp int
}

func PrintProgram(p *Program) string {	
	return PrintExpr(p.Expr)
}

func PrintExpr(e *Expr) string {
	switch {
	case e.Num != nil:
		return string(e.Num.Value)
	case e.Var != nil:
		if e.Var.Name != "" {
			return e.Var.Name
		} else {
			return "tmp" + string(e.Var.Temp)
		}
	case e.Let != nil:

	default:
		return ""
	}
	return ""
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
