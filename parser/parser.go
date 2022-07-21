package parser

import (
	"fmt"

	"github.com/bmoxb/ikou/tokens"
)

func Parse(toks []tokens.Token) (Node, error) {
	p := parser{
		path:          []*Node{},
		nextQuotation: Unquoted,
		lastPosition:  tokens.Position{},
		lastLine:      "",
		lastPopped:    nil,
	}

	for _, tok := range toks {
		err := p.processToken(tok)

		if err != nil {
			return Node{}, err
		}
	}

	n := len(p.path)

	if n > 0 {
		plural := ""
		if n > 1 {
			plural = "s"
		}

		return Node{}, p.makeError(fmt.Sprintf("Mismatched brackets - expected %d more closing bracket%v", n, plural))
	}

	return *p.lastPopped, nil
}

type parser struct {
	path          []*Node
	nextQuotation Quotation
	lastPosition  tokens.Position
	lastLine      string
	lastPopped    *Node
}

func (p *parser) processToken(tok tokens.Token) error {
	var err error

	p.lastPosition = tok.Position
	p.lastLine = tok.OriginalString // TODO

	switch tok.Type {
	case tokens.Open:
		child := newEmptyListNode(p.nextQuotationAndReset(), tok.Position)
		p.push(child)

	case tokens.Quote:
		err = p.setNextQuotation(Quoted)

	case tokens.Backquote:
		err = p.setNextQuotation(Backquoted)

	case tokens.Close:
		err = p.pop()

	default:
		node := newTokenNode(p.nextQuotationAndReset(), tok)
		p.top().addChild(node)
	}

	return err
}

func (p *parser) top() *Node {
	return p.path[len(p.path)-1]
}

func (p *parser) push(node *Node) {
	if len(p.path) > 0 {
		p.top().addChild(node)
	}
	p.path = append(p.path, node)
}

func (p *parser) pop() error {
	n := len(p.path)

	if len(p.path) < 1 {
		return p.makeError("Mismatched brackets - unexpected closing bracket")
	}

	p.lastPopped = p.top()
	p.path = p.path[:n-1]

	return nil
}

func (p *parser) setNextQuotation(q Quotation) error {
	if p.nextQuotation != Unquoted {
		return p.makeError("Cannot have two adjacent quotation markers")
	}

	p.nextQuotation = q
	return nil
}

func (p *parser) nextQuotationAndReset() Quotation {
	q := p.nextQuotation
	p.nextQuotation = Unquoted
	return q
}

func (p *parser) makeError(msg string) error {
	return &ParsingError{pos: p.lastPosition, line: p.lastLine, msg: msg}
}
