package parser

import (
	"strings"

	"github.com/bmoxb/ikou/tokens"
)

type Node struct {
	Content   Content
	Position  tokens.Position
	Quotation Quotation
}

func (n Node) String() string {
	var b strings.Builder

	switch n.Quotation {
	case Quoted:
		b.WriteRune('\'')
	case Backquoted:
		b.WriteByte(',')
	}

	switch c := n.Content.(type) {
	case ListContent:
		b.WriteRune('(')

		for i, child := range c.children {
			b.WriteString(child.String())

			if i < len(c.children)-1 {
				b.WriteByte(' ')
			}
		}

		b.WriteRune(')')

	case TokenContent:
		b.WriteString(c.token.OriginalString)
	}

	return b.String()
}

func (n Node) IsList() bool {
	_, is := n.Content.(ListContent)
	return is
}

func (n Node) Children() []*Node {
	listContent := n.Content.(ListContent)
	return listContent.children
}

func (n Node) Token() tokens.Token {
	tokenContent := n.Content.(TokenContent)
	return tokenContent.token
}

func (n *Node) addChild(child *Node) *Node {
	oldChildren := n.Content.(ListContent).children
	n.Content = ListContent{children: append(oldChildren, child)}
	newChildren := n.Content.(ListContent).children
	return newChildren[len(newChildren)-1]
}

type Content interface {
	content()
}

type ListContent struct {
	children []*Node
}

func (c ListContent) content() {}

type TokenContent struct {
	token tokens.Token
}

func (c TokenContent) content() {}

type Quotation int

const (
	Unquoted Quotation = iota
	Quoted
	Backquoted
)

func newTokenNode(quote Quotation, token tokens.Token) *Node {
	return &Node{
		Content:   TokenContent{token},
		Position:  token.Position,
		Quotation: quote,
	}
}

func newListNode(quote Quotation, children []*Node, pos tokens.Position) *Node {
	return &Node{
		Content:   ListContent{children},
		Position:  pos,
		Quotation: quote,
	}
}

func newEmptyListNode(quote Quotation, pos tokens.Position) *Node {
	return newListNode(quote, []*Node{}, pos)
}
