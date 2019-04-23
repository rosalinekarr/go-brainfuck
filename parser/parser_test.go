package parser

import (
	"bytes"
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/rosalinekarr/go-brainfuck/expr"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("successfully parses pointer increments", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte(">"))
		expected := expr.NewIncrementPtrExpr()

		parser.Parse(reader)

		if parser.ast[0] != expected {
			t.Errorf("expected %T, got: %T", expected, parser.ast[0])
		}
	})

	t.Run("successfully parses pointer decrements", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte("<"))
		expected := expr.NewDecrementPtrExpr()

		parser.Parse(reader)

		if parser.ast[0] != expected {
			t.Errorf("expected %T, got: %T", expected, parser.ast[0])
		}
	})

	t.Run("successfully parses increments", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte("+"))
		expected := expr.NewIncrementExpr()

		parser.Parse(reader)

		if parser.ast[0] != expected {
			t.Errorf("expected %T, got: %T", expected, parser.ast[0])
		}
	})

	t.Run("successfully parses decrements", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte("-"))
		expected := expr.NewDecrementExpr()

		parser.Parse(reader)

		if parser.ast[0] != expected {
			t.Errorf("expected %T, got: %T", expected, parser.ast[0])
		}
	})

	t.Run("successfully parses reads", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte(","))
		expected := expr.NewReadExpr()

		parser.Parse(reader)

		if parser.ast[0] != expected {
			t.Errorf("expected %T, got: %T", expected, parser.ast[0])
		}
	})

	t.Run("successfully parses writes", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte("."))
		expected := expr.NewWriteExpr()

		parser.Parse(reader)

		if parser.ast[0] != expected {
			t.Errorf("expected %T, got: %T", expected, parser.ast[0])
		}
	})

	t.Run("successfully parses loops", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte("[]"))
		expected := expr.NewLoopExpr(nil)

		parser.Parse(reader)

		if reflect.DeepEqual(parser.ast[0], *expected) {
			t.Errorf("expected %T, got: %T", expected, parser.ast[0])
		}
	})

	t.Run("successfully parses nested loops", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte("[[]]"))
		expected := expr.NewLoopExpr([]expr.Expr{
			expr.NewLoopExpr(nil),
		})

		parser.Parse(reader)

		if reflect.DeepEqual(parser.ast[0], *expected) {
			t.Errorf("expected %T, got: %T", expected, parser.ast[0])
		}
	})

	t.Run("returns ErrUnexpectedEOF on unclosed loops", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte("[[]"))
		expected := expr.NewLoopExpr([]expr.Expr{
			expr.NewLoopExpr(nil),
		})

		err := parser.Parse(reader)

		if err != ErrUnexpectedEOF {
			t.Errorf("expected %T, got: %v", expected, err)
		}
	})

	t.Run("returns ErrUnexpectedEOF on unmatched loop closings", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte("[]]"))
		expected := expr.NewLoopExpr([]expr.Expr{
			expr.NewLoopExpr(nil),
		})

		err := parser.Parse(reader)

		if err != ErrUnexpectedLoopClose {
			t.Errorf("expected %T, got: %v", expected, err)
		}
	})

	t.Run("ignores non-command characters", func(t *testing.T) {
		t.Parallel()

		parser := &Parser{}
		reader := bytes.NewReader([]byte("abc+def"))
		expected := expr.NewIncrementExpr()

		parser.Parse(reader)

		if parser.ast[0] != expected {
			t.Errorf("expected %T, got: %T", expected, parser.ast[0])
		}
	})
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("successfully prints \"Hello, world!\"", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		parser := NewParser()
		reader := bytes.NewReader([]byte("-[------->+<]>-.-[->+++++<]>++.+++++++..+++.[->+++++<]>+.------------.--[->++++<]>-.--------.+++.------.--------.-[--->+<]>."))

		var out bytes.Buffer
		parser.Parse(reader)
		err := parser.Run(ctx, nil, &out)

		if err != nil {
			t.Errorf("error: %s", err.Error())
		}

		if out.String() != "Hello, world!" {
			t.Errorf("expected \"Hello, world!\", got: %s", out.String())
		}
	})

	t.Run("cancels correctly", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())

		parser := NewParser()
		reader := bytes.NewReader([]byte("+[->+]"))
		parser.Parse(reader)

		var wg sync.WaitGroup
		var err error
		wg.Add(1)
		go func() {
			err = parser.Run(ctx, nil, nil)
			wg.Done()
		}()
		cancel()
		wg.Wait()

		if err == nil || err.Error() != "context canceled" {
			t.Errorf("expected context cancelled error, got: %v", err)
		}
	})
}
