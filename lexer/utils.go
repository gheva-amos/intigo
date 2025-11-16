package lexer

import (
	"unicode"
)

func (l *Lexer) is_end_char(char rune) bool {
	if unicode.IsSpace(char) {
		return true
	}
	for _, r := range l.end_chars {
		if r == char {
			return true
		}
	}
	return false
}

func (l *Lexer) new_token(val any, tp TokenType) *Token {
	return NewToken(l.Line(), l.Column(), val, tp)
}

func (l *Lexer) handle_double_char(current rune, ret *Token, dcd double_char_data) {
	ret.Type = dcd.IfSingle
	r, eof := l.Peek()
	if eof {
		return
	}
	if r == dcd.DoubleChar {
		ret.Type = dcd.IfDouble
		ret.Value = string(current) + string(r)
		l.Next() // consume the second char
	}
}
