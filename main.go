package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rosalinekarr/go-brainfuck/parser"
)

func main() {
	ctx := context.Background()

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %s\n", err.Error())
		os.Exit(1)
	}

	p := parser.NewParser()
	if err = p.Parse(file); err != nil {
		fmt.Fprintf(os.Stderr, "parsing error: %s\n", err.Error())
		os.Exit(1)
	}

	if err = p.Run(ctx, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %s\n", err.Error())
		os.Exit(1)
	}
}
