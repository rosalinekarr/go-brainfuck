package expr

import (
	"bufio"
)

const MEMSIZE = 256

type Context struct {
	Mem []byte
	Ptr byte
	In  *bufio.Reader
	Out *bufio.Writer
}

func NewContext(in *bufio.Reader, out *bufio.Writer) *Context {
	return &Context{
		Mem: make([]byte, MEMSIZE),
		Ptr: 0,
		In:  in,
		Out: out,
	}
}

type Expr interface {
	Execute(context *Context) error
}

type IncrementPtrExpr struct{}

func NewIncrementPtrExpr() *IncrementPtrExpr {
	return &IncrementPtrExpr{}
}

func (expr *IncrementPtrExpr) Execute(context *Context) error {
	context.Ptr++
	return nil
}

type DecrementPtrExpr struct{}

func NewDecrementPtrExpr() *DecrementPtrExpr {
	return &DecrementPtrExpr{}
}

func (expr *DecrementPtrExpr) Execute(context *Context) error {
	context.Ptr--
	return nil
}

type IncrementExpr struct{}

func NewIncrementExpr() *IncrementExpr {
	return &IncrementExpr{}
}

func (expr *IncrementExpr) Execute(context *Context) error {
	context.Mem[context.Ptr]++
	return nil
}

type DecrementExpr struct{}

func NewDecrementExpr() *DecrementExpr {
	return &DecrementExpr{}
}

func (expr *DecrementExpr) Execute(context *Context) error {
	context.Mem[context.Ptr]--
	return nil
}

type ReadExpr struct{}

func NewReadExpr() *ReadExpr {
	return &ReadExpr{}
}

func (expr *ReadExpr) Execute(context *Context) error {
	inByte, err := context.In.ReadByte()
	if err != nil {
		return err
	}
	context.Mem[context.Ptr] = inByte
	return nil
}

type WriteExpr struct{}

func NewWriteExpr() *WriteExpr {
	return &WriteExpr{}
}

func (expr *WriteExpr) Execute(context *Context) error {
	defer context.Out.Flush()
	outByte := context.Mem[context.Ptr]
	return context.Out.WriteByte(outByte)
}

type LoopExpr struct {
	ast []Expr
}

func NewLoopExpr(ast []Expr) *LoopExpr {
	return &LoopExpr{ast: ast}
}

func (expr *LoopExpr) Execute(context *Context) error {
	for context.Mem[context.Ptr] != 0 {
		for _, subExpr := range expr.ast {
			if err := subExpr.Execute(context); err != nil {
				return err
			}
		}
	}
	return nil
}
