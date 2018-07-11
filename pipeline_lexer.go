package main

type tokenType int

const (
	identifier_token tokenType = iota
	string_literal_token
	pipe_token
	plus_token
	minus_token
)

type token struct {
	tokenType  tokenType
	tokenValue string
}

func lex(cmd string) []token {
	tokens := []token{}
	var readingLiteral bool
	literal := ""

	endLiteral := func() {
		if literal != "" {
			tokens = append(tokens, token{tokenType: identifier_token, tokenValue: literal})
			literal = ""
		}
	}

	escapeNext := false

	for _, c := range cmd {
		if escapeNext {
			literal += string(c)
			escapeNext = false
		} else {
			switch c {
			case '\\':
				escapeNext = true
			case '"':
				if readingLiteral {
					tokens = append(tokens, token{tokenType: string_literal_token, tokenValue: literal})
				}
				readingLiteral = !readingLiteral
				literal = ""
			case '+':
				endLiteral()
				tokens = append(tokens, token{tokenType: plus_token})
			case '-':
				endLiteral()
				tokens = append(tokens, token{tokenType: minus_token})
			case '|':
				endLiteral()
				tokens = append(tokens, token{tokenType: pipe_token})
			case ' ':
				if readingLiteral {
					literal += string(c)
				} else {
					endLiteral()
				}
			default:
				literal += string(c)
			}
		}
	}

	if readingLiteral {
		panic("unterminated string literal")
	}

	if literal != "" {
		tokens = append(tokens, token{tokenType: identifier_token, tokenValue: literal})
	}

	return tokens
}
