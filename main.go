package main

import (
	"fmt"

	"github.com/ironfang-ltd/ironscript/lexer"
	"github.com/ironfang-ltd/ironscript/parser"
)

func main() {

	script := `
fn add(a, b) {
	return a + b;
}

let result = add(5, 10);

print (result);
`

	l := lexer.New(script)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(program.Debug())
	}
}
