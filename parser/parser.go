package parser

import (
  "bufio"
  "bytes"
  "io"

  "github.com/rosalinekarr/go-brainfuck/expr"
)

type Parser struct {
  ast []expr.Expr
}

func NewParser() *Parser {
  return &Parser{}
}

func (parser *Parser) Parse(reader io.Reader) error {
  ast, err := parse(reader)
  if err != nil {
    return err
  }
  parser.ast = ast
  return nil
}

func parse(reader io.Reader) ([]expr.Expr, error) {
  in := bufio.NewReader(reader)
  var ast []expr.Expr
  c, err := in.ReadByte()
  for err != io.EOF {
    if err != nil {
      return nil, err
    }
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
      substr, err := in.ReadBytes(']')
      if err != nil {
        return nil, err
      }
      byteReader := bytes.NewReader(substr)
      ast, err := parse(byteReader)
      if err != nil {
        return nil, err
      }
      expression = expr.NewLoopExpr(ast)
    }
    if expression != nil {
      ast = append(ast, expression)
    }
    c, err = in.ReadByte()
  }
  return ast, nil
}

func (parser *Parser) Run(reader io.Reader, writer io.Writer) error {
  in := bufio.NewReader(reader)
  out := bufio.NewWriter(writer)
  context := expr.NewContext(in, out)
  for _, expression := range parser.ast {
    err := expression.Execute(context)
    if err != nil {
      return err
    }
  }
  return nil
}
