package main

import (
	"reflect"
	"testing"
)

func runParserTest(t *testing.T, pipeline string, expected []PipelineNode) {
	nodes := parsePipeline(pipeline)
	if !reflect.DeepEqual(nodes, expected) {
		t.Logf("FAIL %s", pipeline)
		t.Fatalf("Expected: %+v, Found: %+v", expected, nodes)
	} else {
		t.Logf("PASS %s", pipeline)
	}
}

func TestPipelineParser(t *testing.T) {
	runParserTest(t, "a", []PipelineNode{
		Command{Command: "a"},
	})
	runParserTest(t, "ab", []PipelineNode{
		Command{Command: "ab"},
	})
	runParserTest(t, "a|b", []PipelineNode{
		Command{Command: "a"},
		Command{Command: "b"},
	})
}

func TestPipelineParserExpressions(t *testing.T) {
	runParserTest(t, "a+b", []PipelineNode{
		UnionNode{
			Left:  Command{Command: "a"},
			Right: Command{Command: "b"},
		},
	})
	runParserTest(t, "a-b", []PipelineNode{
		DifferenceNode{
			Left:  Command{Command: "a"},
			Right: Command{Command: "b"},
		},
	})
	runParserTest(t, "a-b+c", []PipelineNode{
		DifferenceNode{
			Left: Command{Command: "a"},
			Right: UnionNode{
				Left:  Command{Command: "b"},
				Right: Command{Command: "c"},
			},
		},
	})
}
