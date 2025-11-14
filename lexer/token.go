package lexer

import (
	"fmt"
)

type Token struct {
	Line   uint64
	Column uint64
	Value  any
	Type   TokenType
}

func NewToken(line uint64, col uint64, val any, tp TokenType) *Token {
	return &Token{Line: line, Column: col, Value: val, Type: tp}
}

func (t Token) String() string {
	return fmt.Sprintf("%d:%d - %v", t.Line, t.Column, t.Value)
}
