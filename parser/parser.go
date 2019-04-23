package parser

import (
	"bufio"
	"context"
	"errors"
	"io"

	"github.com/rosalinekarr/go-brainfuck/expr"
)

var (
	ErrUnexpectedEOF       = errors.New("unexpected end of file")
	ErrUnexpectedLoopClose = errors.New("unexpected ]")
)

type Parser struct {
	ast []expr.Expr
}

func NewParser() *Parser {
	return &Parser{}
}

func (parser *Parser) Parse(reader io.Reader) error {
	ast, err := parse(reader)
	if err != io.EOF {
		return err
	}
	parser.ast = ast
	return nil
}

func parse(reader io.Reader) ([]expr.Expr, error) {
	in := bufio.NewReader(reader)
	var ast []expr.Expr
	c, err := in.ReadByte()
	for err == nil {
		var expression expr.Expr
		switch c {
		case '>':
			expression = expr.NewIncrementPtrExpr()
		case '<':
			expression = expr.NewDecrementPtrExpr()
		case '+':
			expression = expr.NewIncrementExpr()
		case '-':
			expression = expr.NewDecrementExpr()
		case '.':
			expression = expr.NewWriteExpr()
		case ',':
			expression = expr.NewReadExpr()
		case '[':
			ast, err := parse(in)
			if err == io.EOF {
				return nil, ErrUnexpectedEOF
			}
			if err != ErrUnexpectedLoopClose {
				return nil, err
			}
			expression = expr.NewLoopExpr(ast)
		case ']':
			return ast, ErrUnexpectedLoopClose
		}
		if expression != nil {
			ast = append(ast, expression)
		}
		c, err = in.ReadByte()
	}
	return ast, err
}

func (parser *Parser) Run(ctx context.Context, reader io.Reader, writer io.Writer) error {
	in := bufio.NewReader(reader)
	out := bufio.NewWriter(writer)
	context := expr.NewContext(in, out)
	for _, expression := range parser.ast {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := expression.Execute(context)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
