package main

import ()

// Any old piece of a pipeline.
type PipelineNode interface{}

// A command to run in the pipeline
type Command struct {
	// The command name to run.
	Command string

	// The arguments to give it.
	Arguments []string
}

// A set union operation (+)
type UnionNode struct {
	Left  PipelineNode
	Right PipelineNode
}

// A set difference operation (-)
type DifferenceNode struct {
	Left  PipelineNode
	Right PipelineNode
}

// Parse a given command at the given index of the token stream.
// Return the command node, and the new position in the token stream.
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

// Parse a set union operation (+) at the given index in the token stream.
// The left hand side has already been parsed.
// Return the union node, and the new position in the token stream.
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

// Parse a set difference operation (-) at the given index in the token stream.
// The left hand side has already been parsed.
// Return the difference node, and the new position in the token stream.
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

// Parse a node at the given index in the token stream.
// Return the node, and the new position in the token stream.
// Note that this may recurse (e.g. if an operator is found)
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

// Parse a string pipeline into a slice of PipelineNodes.
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
