package lexer

import (
	"fmt"
	"strings"
)

func (l *Lexer) NextString() (string, error) {
	current, eof := l.Next()
	var sb strings.Builder
	for current != '"' && current != '\'' {
		if current == '\\' {
			current, eof = l.Next()
			switch current {
			case 'n':
				current = '\n'
			case 't':
				current = '\t'
			}
		}
		if eof {
			return "", fmt.Errorf("unexpected end of file in the middle of a string")
		}
		sb.WriteRune(current)
		current, eof = l.Next()
	}
	return sb.String(), nil
}
