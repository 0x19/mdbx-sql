package parser

import "fmt"

type Node interface{}

type SelectStatement struct {
	Columns   []string
	TableName string
	Condition string
}

type Parser struct {
	lexer     *Lexer
	curToken  *Token
	peekToken *Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{}
	p.Init(l)
	return p
}

func (p *Parser) Init(l *Lexer) {
	p.lexer = l
	p.nextToken()
	p.nextToken()
}

func (p *Parser) nextToken() {
	if p.curToken != nil {
		p.curToken.reset()
		tokenPool.Put(p.curToken)
	}
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Parse() (Node, error) {
	switch p.curToken.Type {
	case SELECT:
		return p.parseSelectStatement()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.curToken.Literal)
	}
}

func (p *Parser) parseSelectStatement() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	p.nextToken()
	stmt.Columns = p.parseColumnList()

	if p.curToken.Type != FROM {
		return nil, fmt.Errorf("expected FROM, got %s", p.curToken.Literal)
	}
	p.nextToken()

	if p.curToken.Type != IDENT {
		return nil, fmt.Errorf("expected table name, got %s", p.curToken.Literal)
	}
	stmt.TableName = p.curToken.Literal

	p.nextToken()

	if p.curToken.Type == WHERE {
		p.nextToken()
		stmt.Condition = p.curToken.Literal
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseColumnList() []string {
	columns := make([]string, 0)
	columns = append(columns, p.curToken.Literal)

	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		columns = append(columns, p.curToken.Literal)
	}

	p.nextToken()
	return columns
}
