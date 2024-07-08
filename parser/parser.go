package parser

// AST Nodes
type Node interface{}

type SelectStatement struct {
	Columns   []string
	TableName string
	Condition string
}

// Parser represents a parser for SQL.
type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
}

// NewParser initializes a new Parser.
func NewParser(l *Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Parse() Node {
	switch p.curToken.Type {
	case SELECT:
		return p.parseSelectStatement()
	default:
		return nil
	}
}

func (p *Parser) parseSelectStatement() *SelectStatement {
	stmt := &SelectStatement{}

	p.nextToken()
	stmt.Columns = p.parseColumnList()

	if p.curToken.Type != FROM {
		return nil
	}
	p.nextToken()

	if p.curToken.Type != IDENT {
		return nil
	}
	stmt.TableName = p.curToken.Literal

	p.nextToken()

	if p.curToken.Type == WHERE {
		p.nextToken()
		stmt.Condition = p.curToken.Literal
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseColumnList() []string {
	columns := []string{}
	columns = append(columns, p.curToken.Literal)

	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		columns = append(columns, p.curToken.Literal)
	}

	p.nextToken()
	return columns
}
