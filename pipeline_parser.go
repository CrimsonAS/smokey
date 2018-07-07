package main

// A command represents a piece of a pipeline.
type Command struct {
	// The command name to run.
	Command string

	// The arguments to give it.
	Arguments []string
}

type tokenType int

const (
	identifier_token tokenType = iota
	string_literal_token
	pipe_token
)

type token struct {
	tokenType  tokenType
	tokenValue string
}

func lex(cmd string) []token {
	tokens := []token{}
	var readingLiteral bool
	literal := ""

	for _, c := range cmd {
		switch c {
		case '"':
			if readingLiteral {
				tokens = append(tokens, token{tokenType: string_literal_token, tokenValue: literal})
			}
			readingLiteral = !readingLiteral
			literal = ""
		case '|':
			tokens = append(tokens, token{tokenType: pipe_token})
		case ' ':
			if readingLiteral {
				literal += string(c)
			} else if literal != "" {
				tokens = append(tokens, token{tokenType: identifier_token, tokenValue: literal})
				literal = ""
			}
		default:
			literal += string(c)
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

// Parse a pipeline into a slice of command instances.
// For instance, echo foo | cat will turn into two commands: echo (foo), and cat.
func parsePipeline(cmd string) []Command {
	tokens := lex(cmd)
	cmds := []Command{}

	for idx := 0; idx < len(tokens); idx++ {
		tok := tokens[idx]
		switch tok.tokenType {
		case pipe_token:
		case string_literal_token:
			panic("expected: command")
		}
		cmd := Command{Command: tok.tokenValue}

		for idx = idx + 1; /* skip command */ idx < len(tokens) && tokens[idx].tokenType != pipe_token; idx++ {
			switch tok.tokenType {
			case identifier_token:
			case string_literal_token:
			default:
				panic("expected: identifier or string literal")
			}
			cmd.Arguments = append(cmd.Arguments, tokens[idx].tokenValue)
		}

		cmds = append(cmds, cmd)
	}

	return cmds
}
