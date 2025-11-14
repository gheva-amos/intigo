package lexer_test

import (
	"github.com/gheva-amos/intigo/config"
	"github.com/gheva-amos/intigo/lexer"
	"testing"
)

func TestLexerConstructor(t *testing.T) {
	lexer := lexer.New()

	if lexer.Line() != 1 {
		t.Errorf("token.Line got %d want %d", lexer.Line(), 1)
	}
}

func TestLexerDefinition(t *testing.T) {
	cfg := config.ConfigFromJson([]byte(configuration))
	if cfg == nil {
		t.Errorf("Could not instantiate the config")
	}
	lexer := lexer.DefineLexer(cfg)
	if lexer.Line() != 1 {
		t.Errorf("token.Line got %d want %d", lexer.Line(), 1)
	}
}

func TestLexerNextWord(t *testing.T) {
	cfg := config.ConfigFromJson([]byte(configuration))
	if cfg == nil {
		t.Errorf("Could not instantiate the config")
	}
	lexer := lexer.DefineLexer(cfg)

	src := "     test    test2 "
	lexer.AddSource(src)

	word := lexer.NextWord()
	if word != "test" {
		t.Errorf("exepected \"test\" found %s", word)
	}
	word = lexer.NextWord()
	if word != "test2" {
		t.Errorf("exepected \"test2\" found %s", word)
	}
}

const configuration = `{
	"lexer": {
		"token_types": [
			"Integer",
			"Float",
			"Boolean",
			"Plus",
			"Minus",
			"Times",
			"Div",
			"Comma",
			"Equal",
			"EqualEqual",
			"Not",
			"NotEqual",
			"LBrace",
			"RBrace",
			"LParen",
			"RParen",
			"If",
			"Else",
			"For",
			"True",
			"False",
			"Func",
			"Identifier",
			"Return"
		],
		"keyword_map": {
			"if": "If",
			"else": "Else",
			"for": "For",
			"true": "True",
			"false": "False",
			"func": "Func",
			"return": "Return"
		},
		"single_chars": {
			"+": "Plus",
			"-": "Minus",
			"*": "Times",
			"/": "Div",
			"{": "LBrace",
			"}": "RBrace",
			"(": "LParen",
			")": "RParen",
			",": "Comma"
		},
		"double_chars": [
			{
				"if_single": "=",
				"single_type": "Equal",
				"if_double": "=",
				"double_type": "EqualEqual"
			},
		  {
				"if_single": "!",
				"single_type": "Not",
				"if_double": "=",
				"double_type": "NotEqual"
			}
		]
	}
}`
