package parser

import (
	"fmt"
	"monkey/internal/ast"
	"monkey/internal/lexer"
	"monkey/internal/token"
	"strconv"
)

type (
	// 前缀解析函数
	prefixParseFn func() ast.Expression
	// 中缀解析函数
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	// 低优先级
	LOWEST
	// ==
	EQUALS
	// 比较
	LESSGREATER
	// +
	SUM
	// *
	PRODUCT
	// 前缀表达式
	PREFIX
	// 括号、函数调用
	CALL
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l,
		errors: []string{},
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}
	// 读取两个token，设置curToken和peekToken
	p.nextToken()
	p.nextToken()
	// 注册前缀解析函数
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := "expected next token to be %s, got %s instread"
	p.errors = append(p.errors, msg)
}

// 期望下一个token是t类型
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		// 读取下一个token
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		// 解析语句
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		// 读取下一个token
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		// 解析let语句
		return p.parseLetStatement()
	case token.RETURN:
		// 解析return语句
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement();
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {

	stmt := &ast.LetStatement{Token: p.curToken}

	// 解析let后面的标识符
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// 标识符
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 解析等号
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 解析表达式语句
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	// 解析表达式
	stmt.Expression = p.parseExpression(LOWEST)
	// 解析分号
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// 解析表达式
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// 前缀解析函数
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	// 解析前缀表达式
	leftExp := prefix()
	return leftExp
}

// 解析标识符
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// 解析整型字面量
func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.curToken}
	// 解析整型字面量
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := "could not parse %q as integer"
		p.errors = append(p.errors, fmt.Sprintf(msg, p.curToken.Literal))
		return nil
	}
	// 整型字面量值
	literal.Value = value
	return literal
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := "no prefix parse function for %s found"
	p.errors = append(p.errors, fmt.Sprintf(msg, t))
}

// 解析前缀表达式
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	// 读取下一个token
	p.nextToken()
	// 解析右侧表达式
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

