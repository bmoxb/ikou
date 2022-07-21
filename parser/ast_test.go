package parser

import (
	"testing"

	"github.com/bmoxb/ikou/tokens"
)

func TestAddingChildren(t *testing.T) {
	n := newEmptyListNode(Unquoted, tokens.Position{})
	n.addChild(newTokenNode(Unquoted, tokens.Token{}))
	n.addChild(newEmptyListNode(Quoted, tokens.Position{}))

	if len(n.Children()) != 2 {
		t.Errorf("Children not successfully added to empty list as expected length of 2 but found length: %d", len(n.Children()))
	}
}

func TestNodeString(t *testing.T) {
	x := newListNode(
		Unquoted,
		[]*Node{
			newListNode(
				Quoted,
				[]*Node{
					newTokenNode(Unquoted, tokens.Token{OriginalString: "12.5"}),
					newTokenNode(Backquoted, tokens.Token{OriginalString: "abc"}),
				},
				tokens.Position{},
			),
			newTokenNode(Unquoted, tokens.Token{OriginalString: `"hello"`}),
		},
		tokens.Position{},
	)
	table := map[*Node]string{
		x: `('(12.5 ,abc) "hello")`,
	}

	for node, expected := range table {
		if node.String() != expected {
			t.Errorf("Syntax tree incorrectly converted to text as expected %v does not equal encountered %v", expected, node)
		}
	}
}
