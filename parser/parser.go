package parser

// AST Nodes
type Node interface{}

type SelectStatement struct {
	Columns    []string
	TableName  string
	Conditions []Condition
	Joins      []Join
	Aggregates []Aggregate
}

type Condition struct {
	Field string
	Op    string
	Value string
}

type Join struct {
	Table    string
	OnField1 string
	OnField2 string
}

type Aggregate struct {
	Func   string
	Column string
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
	//log.Printf("Parse: curToken=%v", p.curToken)
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
	stmt.Columns, stmt.Aggregates = p.parseColumnList()

	if p.curToken.Type != FROM {
		//log.Printf("parseSelectStatement: expected FROM, got %v", p.curToken)
		return nil
	}
	p.nextToken()

	if p.curToken.Type != IDENT {
		//log.Printf("parseSelectStatement: expected IDENT, got %v", p.curToken)
		return nil
	}
	stmt.TableName = p.curToken.Literal

	p.nextToken()

	for p.curToken.Type == JOIN {
		p.nextToken()
		join := p.parseJoin()
		stmt.Joins = append(stmt.Joins, join)
	}

	if p.curToken.Type == WHERE {
		p.nextToken()
		stmt.Conditions = p.parseConditions()
	}

	return stmt
}

func (p *Parser) parseColumnList() ([]string, []Aggregate) {
	columns := []string{}
	aggregates := []Aggregate{}

	for p.curToken.Type == IDENT || p.curToken.Type == SUM {
		if p.curToken.Type == IDENT {
			columns = append(columns, p.curToken.Literal)
			p.nextToken()
		} else if p.curToken.Type == SUM {
			aggregate := p.parseAggregate()
			aggregates = append(aggregates, aggregate)
		}

		if p.curToken.Type == COMMA {
			p.nextToken()
		} else {
			break
		}
	}

	return columns, aggregates
}

func (p *Parser) parseConditions() []Condition {
	conditions := make([]Condition, 0)

	for {
		cond := Condition{}

		if p.curToken.Type != IDENT {
			//log.Printf("parseConditions: expected IDENT, got %v", p.curToken)
			return nil
		}
		cond.Field = p.curToken.Literal

		p.nextToken()

		if p.curToken.Type != EQ && p.curToken.Type != GT && p.curToken.Type != LT &&
			p.curToken.Type != GTE && p.curToken.Type != LTE && p.curToken.Type != NEQ {
			//log.Printf("parseConditions: expected comparison operator, got %v", p.curToken)
			return nil
		}
		cond.Op = p.curToken.Literal

		p.nextToken()

		if p.curToken.Type != IDENT && p.curToken.Type != NUMBER {
			//log.Printf("parseConditions: expected IDENT or NUMBER, got %v", p.curToken)
			return nil
		}
		cond.Value = p.curToken.Literal

		conditions = append(conditions, cond)

		p.nextToken()

		if p.curToken.Type != AND {
			break
		}

		p.nextToken()
	}

	return conditions
}

func (p *Parser) parseJoin() Join {
	var join Join

	if p.curToken.Type != IDENT {
		//log.Printf("parseJoin: expected IDENT, got %v", p.curToken)
		return join
	}
	join.Table = p.curToken.Literal

	p.nextToken()

	if p.curToken.Type != ON {
		//log.Printf("parseJoin: expected ON, got %v", p.curToken)
		return join
	}
	p.nextToken()

	if p.curToken.Type != IDENT {
		//log.Printf("parseJoin: expected IDENT, got %v", p.curToken)
		return join
	}
	join.OnField1 = p.curToken.Literal

	p.nextToken()

	if p.curToken.Type != EQ {
		//log.Printf("parseJoin: expected EQ, got %v", p.curToken)
		return join
	}
	p.nextToken()

	if p.curToken.Type != IDENT {
		//log.Printf("parseJoin: expected IDENT, got %v", p.curToken)
		return join
	}
	join.OnField2 = p.curToken.Literal

	p.nextToken()

	return join
}

func (p *Parser) parseAggregate() Aggregate {
	var aggregate Aggregate

	if p.curToken.Type != SUM {
		//log.Printf("parseAggregate: expected SUM, got %v", p.curToken)
		return aggregate
	}
	aggregate.Func = "SUM"

	p.nextToken()

	if p.curToken.Type != LPAREN {
		//log.Printf("parseAggregate: expected LPAREN, got %v", p.curToken)
		return aggregate
	}
	p.nextToken()

	if p.curToken.Type != IDENT {
		//log.Printf("parseAggregate: expected IDENT, got %v", p.curToken)
		return aggregate
	}
	aggregate.Column = p.curToken.Literal

	p.nextToken()

	if p.curToken.Type != RPAREN {
		//log.Printf("parseAggregate: expected RPAREN, got %v", p.curToken)
		return aggregate
	}
	p.nextToken()

	return aggregate
}
