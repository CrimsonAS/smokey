package main

import ()

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

	for _, c := range cmd {
		switch c {
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

	if readingLiteral {
		panic("unterminated string literal")
	}

	if literal != "" {
		tokens = append(tokens, token{tokenType: identifier_token, tokenValue: literal})
	}

	return tokens
}

type PipelineNode interface{}

type UnionNode struct {
	Left  PipelineNode
	Right PipelineNode
}

type DifferenceNode struct {
	Left  PipelineNode
	Right PipelineNode
}

func parseCommand(tokens []token, idx int) (PipelineNode, int) {
	switch tokens[idx].tokenType {
	case plus_token:
		panic("Unexpected plus")
	case minus_token:
		panic("Unexpected minus")
	case pipe_token:
		panic("Unexpected pipe")
	case string_literal_token:
		panic("expected: command")
	}

	cmd := Command{Command: tokens[idx].tokenValue}

	for {
		idx = idx + 1
		if idx >= len(tokens) {
			break
		}

		ct := tokens[idx]
		if ct.tokenType == identifier_token || ct.tokenType == string_literal_token {
			cmd.Arguments = append(cmd.Arguments, ct.tokenValue)
		} else {
			break
		}
	}

	return cmd, idx
}

func parsePlus(tokens []token, idx int, lhs PipelineNode) (PipelineNode, int) {
	if tokens[idx].tokenType != plus_token {
		panic("Unexpected non-plus token")
	}
	idx++

	if idx >= len(tokens) {
		panic("Unexpected lack of RHS")
	}

	rhs, idx := parseRecursively(tokens, idx)
	return UnionNode{Left: lhs, Right: rhs}, idx
}

func parseMinus(tokens []token, idx int, lhs PipelineNode) (PipelineNode, int) {
	if tokens[idx].tokenType != minus_token {
		panic("Unexpected non-minus token")
	}
	idx++

	if idx >= len(tokens) {
		panic("Unexpected lack of RHS")
	}

	rhs, idx := parseRecursively(tokens, idx)
	return DifferenceNode{Left: lhs, Right: rhs}, idx
}

func parseRecursively(tokens []token, idx int) (PipelineNode, int) {
	ret, idx := parseCommand(tokens, idx)

	if idx < len(tokens) {
		switch tokens[idx].tokenType {
		case plus_token:
			ret, idx = parsePlus(tokens, idx, ret)
			return ret, idx
		case minus_token:
			ret, idx = parseMinus(tokens, idx, ret)
			return ret, idx
		case pipe_token:
			idx += 1
			return ret, idx
		case identifier_token:
			panic("Unexpected identifier_token")
		case string_literal_token:
			panic("Unexpected string_literal_token")
		}
	} else {
		return ret, idx
	}

	panic("Unreachable")
}

// Parse a pipeline into a slice of command instances.
// For instance, echo foo | cat will turn into two commands: echo (foo), and cat.
func parsePipeline(cmd string) []PipelineNode {
	tokens := lex(cmd)
	nodes := []PipelineNode{}

	for idx := 0; idx < len(tokens); {
		var ret PipelineNode
		ret, idx = parseRecursively(tokens, idx)
		nodes = append(nodes, ret)
	}

	return nodes
}
