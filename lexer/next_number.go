package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// TODO scientific notation
func (l *Lexer) NextNumber(current rune) (any, error) {
	var sb strings.Builder

	var is_hex = false
	var is_oct = false
	var is_bin = false
	var is_float = false
	var seen_dot = false

	var eof bool
	if current == '0' {
		current, eof = l.Peek()
		if eof {
			return nil, fmt.Errorf("Unexpected EOF")
		}
		is_oct = true

		switch {
		case strings.ContainsRune("xX", current):
			is_hex = true
			is_oct = false
			sb.WriteRune('0')
			sb.WriteRune('X')
			current, eof = l.Next()
		case strings.ContainsRune("bB", current):
			is_bin = true
			is_oct = false
			sb.WriteRune('0')
			sb.WriteRune('B')
			current, eof = l.Next()
		case strings.ContainsRune("oO", current):
			sb.WriteRune('0')
			sb.WriteRune('O')
			current, eof = l.Next()
		case current == '.':
			is_oct = false
		}
		if eof {
			return nil, fmt.Errorf("Unexpected EOF")
		}
	}
	// we are at the start of the number proper
	for {
		if eof {
			break
		}
		if unicode.IsDigit(current) {
			sb.WriteRune(current)
		} else if is_hex {
			if strings.ContainsRune("abcdefABCDEF", current) {
				sb.WriteRune(current)
			} else {
				return nil, fmt.Errorf("Found unexpected char looking for a number %c", current)
			}
		} else if current == '.' {
			if seen_dot {
				return nil, fmt.Errorf("Found an extra dot")
			}
			if is_bin || is_oct || is_hex {
				return nil, fmt.Errorf("Found an extra dot")
			}
			seen_dot = true
			is_float = true
			sb.WriteRune(current)
		} else if current == 'e' || current == 'E' {
			is_float = true
			sb.WriteRune(current)
		} else if current == '-' && is_float {
			sb.WriteRune(current)
		} else if l.is_end_char(current) {
			break
		}
		current, eof = l.Next()
	}
	if !is_float {
		return strconv.ParseInt(sb.String(), 0, 64)
	}
	return strconv.ParseFloat(sb.String(), 64)
}
