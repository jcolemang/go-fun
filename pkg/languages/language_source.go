package languages

import (
	"fmt"

    "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/alecthomas/repr"
)

// Main Language
// This is the source language that this compiler is targeted for

type Empty struct {}

type Program struct {
	Expr *Expr `@@`
}

type Expr struct {
	Bool *Bool     `@@`
	Num *Num       `| @@`
	Let *LetExpr   `| @@`
    IfExpr *IfExpr `| @@`
    App []*Expr    `| "(" @@ @@* ")" `
    Var *Var       `| @@`
}

type IfExpr struct {
	IfCond *Expr  `"(" "if" @@`
	IfTrue *Expr  `@@`
    IfFalse *Expr `@@ ")"`
}

type LetExpr struct {
	LetAssignments []*Assignment[Expr] `"(" "let" "(" @@ @@* ")"`
	LetBody *Expr                      `@@ ")"`
}

type Assignment [T any] struct {
	Ref *Var   `"(" @@`
	Expr *T `@@ ")"`
}

type Primitive struct {
	Num *Num
	Var *Var
    Bool *Bool
}

type PrimitiveApplication struct {
	Operator *Var
	Operands []*Primitive
}

// have this separated to also handle floats
type Num struct {
	Int *int `@Int`
}

type Bool struct {
	True *string  `@"#True"`
	False *string `| @"#False"`
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

func ProgToString(prog *Program) string {
    return repr.String(prog)
}

func GetLanguageParser() *participle.Parser[Program] {
	basicLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"String", `"(\\"|[^"])*"`},
		{"Ident", `[a-zA-Z_\-+#]\w*`}, // note that allowing # allows variables to attempt but fail to shadow bools.
                                       // could maybe fancy this up with some kind of stateful parsing but for now
                                       // I just simply don't care
		{"Int", `[-]?(\d+)`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{"whitespace", `[ \t\n\r]+`},
	})
    parser := participle.MustBuild[Program](
        participle.Lexer(basicLexer),
    )
    return parser
}
