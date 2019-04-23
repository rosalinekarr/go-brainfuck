package expr

import (
	"bufio"
	"bytes"
	"testing"
)

func TestIncrementPtrExpr(t *testing.T) {
	t.Parallel()

	t.Run("increments the pointer", func(t *testing.T) {
		t.Parallel()

		context := &Context{
			Ptr: 0,
		}

		NewIncrementPtrExpr().Execute(context)

		if context.Ptr != 1 {
			t.Errorf("expected 1, got: %d", context.Ptr)
		}
	})

	t.Run("rolls over at 255", func(t *testing.T) {
		t.Parallel()

		context := &Context{
			Ptr: 255,
		}

		NewIncrementPtrExpr().Execute(context)

		if context.Ptr != 0 {
			t.Errorf("expected 0, got: %d", context.Ptr)
		}
	})
}

func TestDecrementPtrExpr(t *testing.T) {
	t.Parallel()

	t.Run("decrements the pointer", func(t *testing.T) {
		t.Parallel()

		context := &Context{
			Ptr: 1,
		}

		NewDecrementPtrExpr().Execute(context)

		if context.Ptr != 0 {
			t.Errorf("expected 0, got: %d", context.Ptr)
		}
	})

	t.Run("rolls over at 0", func(t *testing.T) {
		t.Parallel()

		context := &Context{
			Ptr: 0,
		}

		NewDecrementPtrExpr().Execute(context)

		if context.Ptr != 255 {
			t.Errorf("expected 255, got: %d", context.Ptr)
		}
	})
}

func TestIncrementExpr(t *testing.T) {
	t.Parallel()

	t.Run("increments memory at the pointer", func(t *testing.T) {
		t.Parallel()

		context := &Context{
			Mem: []byte{0},
			Ptr: 0,
		}

		NewIncrementExpr().Execute(context)

		if context.Mem[0] != 1 {
			t.Errorf("expected 1, got: %d", context.Mem[0])
		}
	})

	t.Run("rolls over at 255", func(t *testing.T) {
		t.Parallel()

		context := &Context{
			Mem: []byte{255},
			Ptr: 0,
		}

		NewIncrementExpr().Execute(context)

		if context.Mem[0] != 0 {
			t.Errorf("expected 0, got: %d", context.Mem[0])
		}
	})
}

func TestDecrementExpr(t *testing.T) {
	t.Parallel()

	t.Run("decrements memory at the pointer", func(t *testing.T) {
		t.Parallel()

		context := &Context{
			Mem: []byte{1},
			Ptr: 0,
		}

		NewDecrementExpr().Execute(context)

		if context.Mem[0] != 0 {
			t.Errorf("expected 0, got: %d", context.Mem[0])
		}
	})

	t.Run("rolls over at 0", func(t *testing.T) {
		t.Parallel()

		context := &Context{
			Mem: []byte{0},
			Ptr: 0,
		}

		NewDecrementExpr().Execute(context)

		if context.Mem[0] != 255 {
			t.Errorf("expected 255, got: %d", context.Mem[0])
		}
	})
}

func TestReadExpr(t *testing.T) {
	t.Parallel()

	t.Run("sets memory at the pointer with a byte read from input", func(t *testing.T) {
		t.Parallel()

		reader := bytes.NewReader([]byte{42})

		context := &Context{
			Mem: []byte{0},
			Ptr: 0,
			In:  bufio.NewReader(reader),
		}

		NewReadExpr().Execute(context)

		if context.Mem[0] != 42 {
			t.Errorf("expected 42, got: %d", context.Mem[0])
		}
	})

	t.Run("does nothing on EOF", func(t *testing.T) {
		t.Parallel()

		reader := bytes.NewReader([]byte{})

		context := &Context{
			Mem: []byte{0},
			Ptr: 0,
			In:  bufio.NewReader(reader),
		}

		NewReadExpr().Execute(context)

		if context.Mem[0] != 0 {
			t.Errorf("expected 0, got: %d", context.Mem[0])
		}
	})
}

func TestWriteExpr(t *testing.T) {
	t.Parallel()

	t.Run("prints the current memory value at the pointer to output", func(t *testing.T) {
		t.Parallel()

		writer := new(bytes.Buffer)

		context := &Context{
			Mem: []byte{42},
			Ptr: 0,
			Out: bufio.NewWriter(writer),
		}

		NewWriteExpr().Execute(context)

		if writer.Bytes()[0] != 42 {
			t.Errorf("expected 42, got: %d", writer.Bytes()[0])
		}
	})
}

func TestLoopExpr(t *testing.T) {
	t.Parallel()

	ast := []Expr{}
	ast = append(ast, NewDecrementExpr())
	ast = append(ast, NewIncrementPtrExpr())
	ast = append(ast, NewIncrementExpr())
	ast = append(ast, NewDecrementPtrExpr())

	context := &Context{
		Mem: []byte{42, 0},
		Ptr: 0,
	}

	NewLoopExpr(ast).Execute(context)

	t.Run("loops until the memory value at the pointer is zero", func(t *testing.T) {
		t.Parallel()

		if context.Mem[0] != 0 {
			t.Errorf("expected 0, got: %d", context.Mem[0])
		}
	})

	t.Run("executes abstract syntax subtree", func(t *testing.T) {
		t.Parallel()

		if context.Mem[1] != 42 {
			t.Errorf("expected 42, got: %d", context.Mem[1])
		}
	})
}
