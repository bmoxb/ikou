package parser

import (
	"reflect"
	"testing"

	"github.com/bmoxb/ikou/tokens"
)

func TestParse(t *testing.T) {
	table := map[*Node][]tokens.Token{
		newEmptyListNode(Unquoted, tokens.Position{}): {},
		newEmptyListNode(Unquoted, tokens.Position{}): {tokens.Token{Type: tokens.Open}, tokens.Token{Type: tokens.Close}},

		newListNode(Unquoted, []*Node{
			newEmptyListNode(Unquoted, tokens.Position{}),
			newEmptyListNode(Unquoted, tokens.Position{}),
		}, tokens.Position{}): {
			tokens.Token{Type: tokens.Open},
			tokens.Token{Type: tokens.Open},
			tokens.Token{Type: tokens.Close},
			tokens.Token{Type: tokens.Open},
			tokens.Token{Type: tokens.Close},
			tokens.Token{Type: tokens.Close},
		},

		newListNode(Quoted, []*Node{
			newEmptyListNode(Backquoted, tokens.Position{}),
			newEmptyListNode(Unquoted, tokens.Position{}),
		}, tokens.Position{}): {
			tokens.Token{Type: tokens.Quote},
			tokens.Token{Type: tokens.Open},
			tokens.Token{Type: tokens.Backquote},
			tokens.Token{Type: tokens.Open},
			tokens.Token{Type: tokens.Close},
			tokens.Token{Type: tokens.Open},
			tokens.Token{Type: tokens.Close},
			tokens.Token{Type: tokens.Close},
		},
	}

	for expected, tokens := range table {
		node, err := Parse(tokens)

		if err != nil {
			t.Errorf("Parse( %v ) - unexpected error: %v", tokens, err)
		} else if !reflect.DeepEqual(node, *expected) {
			t.Errorf("Parse( %v ) - want: %v, got: %v", tokens, *expected, node)
		}
	}
}
