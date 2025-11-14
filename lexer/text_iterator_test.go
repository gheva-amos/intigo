package lexer_test

import (
	"github.com/gheva-amos/intigo/lexer"
	"testing"
)

func TestTextIteratorCreation(t *testing.T) {
	ti := lexer.NewTextIterator()

	if ti.Line() != 1 {
		t.Errorf("token.Line got %d want %d", ti.Line(), 1)
	}
}

func TestTextIteratorSkipWhites(t *testing.T) {
	src1 := "  \t\n1"

	ti := lexer.NewTextIterator()
	ti.AddSource(src1)
	eof := ti.SkipWhites()
	if eof {
		t.Errorf("found EOF expected to find a char")
	}
	r, eof := ti.Next()
	if eof {
		t.Errorf("found EOF expected to find a char")
	}
	if r != '1' {
		t.Errorf("expected to find the char '1', found '%c'", r)
	}
	if ti.Line() != 2 {
		t.Errorf("Expected to be in the second line found %d", ti.Line())
	}
	src2 := "  \t\n    "
	ti.AddSource(src2)
	eof = ti.SkipWhites()
	if !eof {
		t.Errorf("Expect EOF")
	}
	src3 := "no white space at the start"
	ti.AddSource(src3)
	eof = ti.SkipWhites()
	if eof {
		t.Errorf("Unxpected EOF")
	}
	r, eof = ti.Next()
	if eof {
		t.Errorf("found EOF expected to find a char")
	}
	if r != 'n' {
		t.Errorf("expected to find the char 'n', found '%c'", r)
	}
}

func TestTextIteratorNextNonWhite(t *testing.T) {
	src1 := "  \t\n1"

	ti := lexer.NewTextIterator()
	ti.AddSource(src1)
	r, eof := ti.NextNonWhite()
	if eof {
		t.Errorf("found EOF expected to find a char")
	}
	if r != '1' {
		t.Errorf("expected to find the char '1', found '%c'", r)
	}
}
