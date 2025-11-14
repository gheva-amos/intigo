package lexer

import (
	"unicode"
)

func (l *Lexer) has_source() bool {
	return l.Source.HasSource()
}

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
