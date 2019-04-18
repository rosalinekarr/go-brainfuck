package parser

import (
	"bytes"
	"reflect"
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

	t.Run("prints \"Hello, world!\"", func(t *testing.T) {
		t.Parallel()

		parser := NewParser()
		reader := bytes.NewReader([]byte("-[------->+<]>-.-[->+++++<]>++.+++++++..+++.[->+++++<]>+.------------.--[->++++<]>-.--------.+++.------.--------.-[--->+<]>."))

		var out bytes.Buffer
		parser.Parse(reader)
		parser.Run(nil, &out)

		if out.String() != "Hello, world!" {
			t.Errorf("expected \"Hello, world!\", got: %s", out.String())
		}
	})
}
