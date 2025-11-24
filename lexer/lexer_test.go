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

func TestLexerNextNumber(t *testing.T) {
	cfg := config.ConfigFromJson([]byte(configuration))
	if cfg == nil {
		t.Errorf("Could not instantiate the config")
	}
	lexer := lexer.DefineLexer(cfg)

	src := "     1234    1.23 1.23e-4"
	lexer.AddSource(src)

	tests := []struct {
		Expect any
	}{
		{Expect: int64(1234)},
		{Expect: float64(1.23)},
		{Expect: float64(1.23e-4)},
	}
	for _, test := range tests {
		r, eof := lexer.NextNonWhite()
		if eof {
			t.Errorf("Unexpected eof")
		}
		num, err := lexer.NextNumber(r)
		if err != nil {
			t.Errorf("Got an error %s", err)
		}
		if num != test.Expect {
			t.Errorf("Number is %d expecting %v", num, test.Expect)
		}
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

	r, eof := lexer.NextNonWhite()
	if eof {
		t.Errorf("Unexpected eof")
	}
	word := lexer.NextWord(r)
	if word != "test" {
		t.Errorf("exepected \"test\" found %s", word)
	}
	r, eof = lexer.NextNonWhite()
	if eof {
		t.Errorf("Unexpected eof")
	}
	word = lexer.NextWord(r)
	if word != "test2" {
		t.Errorf("exepected \"test2\" found %s", word)
	}
}

func TestLexerNextToken(t *testing.T) {
	cfg := config.ConfigFromJson([]byte(configuration))
	if cfg == nil {
		t.Errorf("Could not instantiate the config")
	}
	lexer := lexer.DefineLexer(cfg)

	src := ` - * 
  + 
  = ! != == 1234 1.23e-4
  " test " 1234
  an_identifier
  for
  `
	lexer.AddSource(src)

	type test_info struct {
		expect any
		tp     string
	}
	tests := []test_info{
		{
			expect: "-",
			tp:     "Minus",
		},
		{
			expect: "*",
			tp:     "Times",
		},
		{
			expect: "+",
			tp:     "Plus",
		},
		{
			expect: "=",
			tp:     "Equal",
		},
		{
			expect: "!",
			tp:     "Not",
		},
		{
			expect: "!=",
			tp:     "NotEqual",
		},
		{
			expect: "==",
			tp:     "EqualEqual",
		},
		{
			expect: int64(1234),
			tp:     "Number",
		},
		{
			expect: 1.23e-4,
			tp:     "Number",
		},
		{
			expect: " test ",
			tp:     "String",
		},
		{
			expect: int64(1234),
			tp:     "Number",
		},
		{
			expect: "an_identifier",
			tp:     "Identifier",
		},
		{
			expect: "for",
			tp:     "For",
		},
	}

	for _, test := range tests {
		token := lexer.NextToken()
		if lexer.TypeName(token.Type) != test.tp {
			t.Errorf("Expected %s, Found %s", test.tp, lexer.TypeName(token.Type))
		}
		if token.Value != test.expect {
			t.Errorf("Expected %v, Found %v", test.expect, token.Value)
		}
	}
}

func TestLexerNextString(t *testing.T) {
	cfg := config.ConfigFromJson([]byte(configuration))
	if cfg == nil {
		t.Errorf("Could not instantiate the config")
	}
	lexer := lexer.DefineLexer(cfg)

	src := ` " test string "  `
	lexer.AddSource(src)

	_, eof := lexer.NextNonWhite()
	if eof {
		t.Errorf("Unexpected eof")
	}
	str, err := lexer.NextString()
	if err != nil {
		t.Errorf("%e", err)
	}
	if str != " test string " {
		t.Errorf("expected ' test string ', found '%s'", str)
	}
}

const configuration = `{
	"lexer": {
		"token_types": [
			"Integer",
			"Float",
			"Boolean",
			"String",
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
