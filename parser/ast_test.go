package parser

import (
	"testing"

	"github.com/bmoxb/ikou/lexer"
)

func TestAddingChildren(t *testing.T) {
	n := newEmptyListNode(Unquoted, lexer.TokenPosition{})
	n.addChild(newTokenNode(Unquoted, lexer.Token{}))
	n.addChild(newEmptyListNode(Quoted, lexer.TokenPosition{}))

	if len(n.Children()) != 2 {
		t.Errorf("Children not successfully added to empty list as expected length of 2 but found length: %d", len(n.Children()))
	}
}

func TestNodeString(t *testing.T) {
	x := newListNode(
		Unquoted,
		[]Node{
			newListNode(
				Quoted,
				[]Node{
					newTokenNode(Unquoted, lexer.Token{OriginalString: "12.5"}),
					newTokenNode(Backquoted, lexer.Token{OriginalString: "abc"}),
				},
				lexer.TokenPosition{},
			),
			newTokenNode(Unquoted, lexer.Token{OriginalString: `"hello"`}),
		},
		lexer.TokenPosition{},
	)
	table := map[*Node]string{
		&x: `('(12.5 ,abc) "hello")`,
	}

	for node, expected := range table {
		if node.String() != expected {
			t.Errorf("Syntax tree incorrectly converted to text as expected %v does not equal encountered %v", expected, node)
		}
	}
}
