package languages

import (
	"fmt"

    "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Main Language
// This is the source language that this compiler is targeted for

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

// have this separated to also handle floats
type Num struct {
	Int *int `@Int`
}

type Var struct {
	Primitive string `@("+" | "print")`
    Name string      `| @Ident`
    Generated int
}

func VarToString(v *Var) string {
	if v.Name != "" {
		return v.Name
	} else if v.Primitive != "" {
		return v.Primitive
	} else {
		return "tmp" + fmt.Sprint(v.Generated)
	}
}

func GetBuiltIns() []*Var {
	return []*Var{
		&Var{Primitive: "+"},
	}
}

func GetVarGenerator() func() *Var {
    current := 0
    generator := func() *Var {
        current++
        return &Var{
			Generated: current,
		}
    }
    return generator
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
