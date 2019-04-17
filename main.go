package main

import (
	"os"

	"github.com/rosalinekarr/go-brainfuck/parser"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	p := parser.NewParser()
	if err = p.Parse(file); err != nil {
		panic(err)
	}

	if err = p.Run(os.Stdin, os.Stdout); err != nil {
		panic(err)
	}
}
