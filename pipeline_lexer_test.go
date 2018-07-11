package main

import (
	"reflect"
	"testing"
)

func runLexerTest(t *testing.T, pipeline string, expected []token) {
	tokens := lex(pipeline)
	if !reflect.DeepEqual(tokens, expected) {
		t.Logf("FAIL %s", pipeline)
		t.Fatalf("Expected: %+v, Found: %+v", expected, tokens)
	} else {
		t.Logf("PASS %s", pipeline)
	}
}

func TestPipelineLexer(t *testing.T) {
	runLexerTest(t, "a", []token{
		token{tokenType: identifier_token, tokenValue: "a"},
	})
	runLexerTest(t, "ab", []token{
		token{tokenType: identifier_token, tokenValue: "ab"},
	})
	runLexerTest(t, "a|b", []token{
		token{tokenType: identifier_token, tokenValue: "a"},
		token{tokenType: pipe_token},
		token{tokenType: identifier_token, tokenValue: "b"},
	})
	runLexerTest(t, "a+b", []token{
		token{tokenType: identifier_token, tokenValue: "a"},
		token{tokenType: plus_token},
		token{tokenType: identifier_token, tokenValue: "b"},
	})
	runLexerTest(t, "a-b", []token{
		token{tokenType: identifier_token, tokenValue: "a"},
		token{tokenType: minus_token},
		token{tokenType: identifier_token, tokenValue: "b"},
	})
	runLexerTest(t, "abcdefg", []token{
		token{tokenType: identifier_token, tokenValue: "abcdefg"},
	})
	runLexerTest(t, "\"a\"", []token{
		token{tokenType: string_literal_token, tokenValue: "a"},
	})
	runLexerTest(t, "\"abcdefg\"", []token{
		token{tokenType: string_literal_token, tokenValue: "abcdefg"},
	})
}
