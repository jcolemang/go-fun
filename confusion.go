package main

import (
    "fmt"
	"github.com/alecthomas/repr"
)

type ISomething interface {
    whatever()
}

type Something struct {
    Value string
}

func (b Something) whatever() {}

func Pointers() {
    many := []ISomething{}

    fmt.Println("I add a value")
    many = append(many, Something{
        Value: "hi",
    })

    fmt.Println("I add a pointer to a value")
    many = append(many, &Something{
        Value: "hi",
    })

    one, two := many[0], many[1]

    fmt.Println("I can access the first")
    repr.Println(one)

    fmt.Println("But I cannot dereference the second")
    // repr.Println(*two)

    switch x := two.(type) {
    case Something:
        fmt.Println("This case does not match")
        repr.Println(x)
    case *Something:
        fmt.Println("This case does match")
        repr.Println(x)
    }
}

type Animal interface {
    creature()
}

type Cat interface {
    meow()
}

type Kitten struct {
    Pounces int
}

func (k Kitten) meow() {}
// func (c Cat) creature() {}

func main() {
    // pretty sure this whole file will make sense to me if I read about method sets
    Pointers()
}
