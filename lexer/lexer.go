// The lexer package takes a json file and creates a tokenizer
// the tokenizer then can tokenize a file written with the correct tokens
package lexer

import (
	"github.com/gheva-amos/intigo/config"
	"strings"
)

type TokenType int

const (
	EOF     TokenType = -666
	Unknown TokenType = -777
)

type Lexer struct {
	Source TextIterator

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

	ret.Source.Reset()
	return ret
}

func DefineLexer(cfg *config.Config) *Lexer {
	ret := New()
	ret.token_types = make(map[string]TokenType)
	ret.token_type_names = make(map[TokenType]string)
	ret.token_type_names[EOF] = "EOF"
	ret.token_type_names[Unknown] = "Unknown"
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

func (l *Lexer) Line() uint64 {
	return l.Source.Line()
}

func (l *Lexer) Column() uint64 {
	return l.Source.Column()
}

func (l *Lexer) AddSource(source string) {
	l.Source.AddSource(source)
}

func (l *Lexer) NextWord() string {
	l.Source.SkipWhites()
	var sb strings.Builder
	for next, eof := l.Source.Next(); !eof && !l.is_end_char(next); {
		sb.WriteRune(next)
		next, eof = l.Source.Next()
	}
	return sb.String()
}
