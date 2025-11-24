// The lexer package takes a json file and creates a tokenizer
// the tokenizer then can tokenize a file written with the correct tokens
package lexer

import (
	"github.com/gheva-amos/intigo/config"
	"strings"
	"unicode"
)

type TokenType int

const (
	EOF     TokenType = -666
	Unknown TokenType = -777
	Number  TokenType = -10
)

type Lexer struct {
	TextIterator

	end_chars []rune

	token_types      map[string]TokenType
	token_type_names map[TokenType]string

	keyword_map map[string]TokenType

	single_chars map[rune]TokenType

	double_chars map[rune]double_char_data
}

type double_char_data struct {
	IfSingle   TokenType
	IfDouble   TokenType
	DoubleChar rune
}

func New() *Lexer {
	ret := &Lexer{}

	ret.Reset()
	return ret
}

func DefineLexer(cfg *config.Config) *Lexer {
	ret := New()
	ret.token_types = make(map[string]TokenType)
	ret.token_type_names = make(map[TokenType]string)
	ret.token_type_names[EOF] = "EOF"
	ret.token_type_names[Unknown] = "Unknown"
	ret.token_type_names[Number] = "Number"
	for k, v := range ret.token_type_names {
		ret.token_types[v] = k
	}
	for i, tp := range cfg.Lexer.TokenTypes {
		ret.token_types[tp] = TokenType(i)
		ret.token_type_names[TokenType(i)] = tp
	}
	ret.keyword_map = make(map[string]TokenType)
	for k, v := range cfg.Lexer.KeywordMap {
		ret.keyword_map[k] = ret.token_types[v]
	}
	seen := make(map[rune]struct{})
	ret.single_chars = make(map[rune]TokenType)
	for k, v := range cfg.Lexer.SingleChars {
		r := []rune(k)[0]
		ret.end_chars = append(ret.end_chars, r)
		ret.single_chars[r] = ret.token_types[v]
		seen[r] = struct{}{}
	}
	ret.double_chars = make(map[rune]double_char_data)
	for _, v := range cfg.Lexer.DoubleChars {
		r := []rune(v.Single)[0]

		if _, exists := seen[r]; !exists {
			ret.end_chars = append(ret.end_chars, r)
			seen[r] = struct{}{}
		}
		ret.double_chars[r] = double_char_data{
			IfSingle:   ret.token_types[v.SingleType],
			IfDouble:   ret.token_types[v.DoubleType],
			DoubleChar: []rune(v.Double)[0],
		}
	}
	return ret
}

func (l *Lexer) TypeName(tp TokenType) string {
	if ret, ok := l.token_type_names[tp]; ok {
		return ret
	}
	return ""
}

func (l *Lexer) NextWord(current rune) string {
	var sb strings.Builder
	sb.WriteRune(current)
	for next, eof := l.Next(); !eof && !l.is_end_char(next); {
		sb.WriteRune(next)
		next, eof = l.Next()
	}
	return sb.String()
}

func (l *Lexer) NextToken() *Token {
	current, eof := l.NextNonWhite()
	if eof {
		return l.new_token(rune(0), l.token_types["EOF"])
	}

	ret := l.new_token(string(current), l.token_types["Unknown"])
	if current == '"' || current == '\'' {
		var err error
		ret.Value, err = l.NextString()
		ret.Type = l.token_types["String"]
		if err != nil {
			return nil
		}
		return ret
	}
	if p, _ := l.Peek(); current == '-' && unicode.IsDigit(p) {
		num, err := l.NextNumber(current)
		if err != nil {
			return nil
		}
		ret.Type = l.token_types["Number"]
		ret.Value = num
		return ret
	}
	if tp, ok := l.single_chars[current]; ok {
		ret.Type = tp
	} else if dc, ok := l.double_chars[current]; ok {
		l.handle_double_char(current, ret, dc)
	} else if unicode.IsDigit(current) {
		num, err := l.NextNumber(current)
		if err != nil {
			return nil
		}
		ret.Type = l.token_types["Number"]
		ret.Value = num
	} else {
		word := l.NextWord(current)
		ret.Value = word
		tp, ok := l.keyword_map[word]
		if ok {
			ret.Type = tp
		} else {
			ret.Type = l.token_types["Identifier"]
		}
	}
	return ret
}
