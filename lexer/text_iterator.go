package lexer

import (
	"unicode"
)

type TextIterator struct {
	line  uint64
	col   uint64
	index uint64
	next  uint64

	source []rune
}

func NewTextIterator() *TextIterator {
	ret := &TextIterator{}
	ret.Reset()
	return ret
}

func (ti *TextIterator) AddSource(src string) {
	ti.source = []rune(src)
	ti.Reset()
}

func (ti *TextIterator) Reset() {
	ti.line = 1
	ti.col = 1
	ti.index = 0
	ti.next = 1
}

func (ti *TextIterator) Line() uint64 {
	return ti.line
}

func (ti *TextIterator) Column() uint64 {
	return ti.col
}

func (ti *TextIterator) Index() uint64 {
	return ti.index
}

func (ti *TextIterator) HasSource() bool {
	return len(ti.source) != 0
}

func (ti *TextIterator) Length() uint64 {
	return uint64(len(ti.source))
}

func (ti *TextIterator) consume() {
	ti.index = ti.next
	ti.next += 1
}

func (ti *TextIterator) eof() bool {
	return ti.index >= uint64(len(ti.source))
}

func (ti *TextIterator) advance() rune {
	ret := ti.source[ti.index]
	ti.consume()

	is_new_line := func(r rune) bool {
		return r == '\n' || r == '\r'
	}

	ti.col += 1
	if is_new_line(ret) {
		ti.col = 1
		ti.line += 1
	}
	return ret
}

// return true if EOF false if a valid rune was returned
func (ti *TextIterator) Peek() (rune, bool) {
	if !ti.eof() {
		return ti.source[ti.index], false
	}
	return 0, true
}

// return true if EOF false if a valid rune was returned
func (ti *TextIterator) SkipWhites() bool {
	for c, _ := ti.Peek(); unicode.IsSpace(c); c, _ = ti.Peek() {
		ti.advance()
	}
	return ti.eof()
}

// return true if EOF false if a valid rune was returned
func (ti *TextIterator) Next() (rune, bool) {
	if ti.eof() {
		return 0, true
	}
	ret := ti.advance()

	return ret, false
}

// return true if EOF false if a valid rune was returned
func (ti *TextIterator) NextNonWhite() (rune, bool) {
	eof := ti.SkipWhites()
	if !eof {
		return ti.Next()
	}
	return 0, eof
}

func (ti *TextIterator) PushCharBack() {
	ti.next = ti.index
	ti.index -= 1
	ti.col -= 1
}
