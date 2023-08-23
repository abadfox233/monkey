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
	infixParseFn func(ast.Expression) ast.Expression
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
	// [
	INDEX
)

// 操作符优先级
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

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
		errors:         []string{},
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
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	// 注册括号解析函数
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// 注册中缀解析函数
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	// 注册函数调用解析函数
	p.registerInfix(token.LPAREN, p.parseCallExpression)
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

func (p *Parser) parseHashLiteral() ast.Expression {

	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		// 解析key
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		// 解析value
		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		// 如果下一个token是逗号，继续解析
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	// 期望下一个token是右括号
	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash

}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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
		return p.parseExpressionStatement()
	}
}

// 解析数组索引表达式
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	// 读取下一个token
	p.nextToken()
	// 解析索引表达式
	exp.Index = p.parseExpression(LOWEST)
	// 期望下一个token是]
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

// 解析数组字面量
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

// 解析表达式列表语句
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	// 如果下一个token是结束符，说明没有表达式
	if p.peekTokenIs(end) {
		// 读取下一个token
		p.nextToken()
		return list
	}
	// 读取下一个token
	p.nextToken()
	// 解析表达式
	list = append(list, p.parseExpression(LOWEST))
	// 如果下一个token是逗号，说明还有表达式
	for p.peekTokenIs(token.COMMA) {
		// 读取下一个token
		p.nextToken()
		// 读取下一个token
		p.nextToken()
		// 解析表达式
		list = append(list, p.parseExpression(LOWEST))
	}
	// 如果下一个token不是结束符，报错
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// 解析函数调用
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	// 解析函数参数
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

// 解析let语句
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
	// 读取下一个token
	p.nextToken()
	// 解析等号右侧的表达式
	stmt.Value = p.parseExpression(LOWEST)
	// 解析分号
	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// 解析函数字面量
func (p *Parser) parseFunctionLiteral() ast.Expression {

	lit := &ast.FunctionLiteral{Token: p.curToken}

	// 解析函数后面的左括号
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	// 解析函数参数
	lit.Parameters = p.parseFunctionParameters()
	// 解析函数后面的左大括号
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	// 解析函数体
	lit.Body = p.parseBlockStatement()
	return lit
}

// 解析 If 表达式
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	// 解析if后面的左括号
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	// 解析if后面的右括号
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	// 解析if后面的左大括号
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		// 解析else后面的左大括号
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()

	}
	return expression
}

// 解析函数参数
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifier := []*ast.Identifier{}
	// 如果下一个token是右括号，说明没有参数
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifier
	}
	// 读取下一个token
	p.nextToken()
	// 解析参数
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifier = append(identifier, ident)
	// 如果下一个token是逗号，说明还有参数
	for p.peekTokenIs(token.COMMA) {
		// 读取下一个token
		p.nextToken()
		// 读取下一个token
		p.nextToken()
		// 解析参数
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifier = append(identifier, ident)
	}
	// 解析右括号
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifier
}

// 解析 BlockStatement
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	// 读取下一个token
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		// 解析语句
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		// 读取下一个token
		p.nextToken()
	}
	return block

}

// 解析return语句
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// 解析return后面的表达式
	stmt.ReturnValue = p.parseExpression(LOWEST)
	// 解析分号
	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 解析表达式语句
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
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
	// defer untrace(trace("parseExpression"))
	// 前缀解析函数
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	// 解析前缀表达式
	leftExp := prefix()

	// 解析中缀表达式
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// 中缀解析函数
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		// 读取下一个token
		p.nextToken()
		// 解析中缀表达式
		leftExp = infix(leftExp)
	}

	return leftExp
}

// 解析括号表达式
func (p *Parser) parseGroupedExpression() ast.Expression {
	// defer untrace(trace("parseGroupedExpression"))
	// 读取下一个token
	p.nextToken()
	// 解析表达式
	exp := p.parseExpression(LOWEST)
	// 解析右括号
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// 解析标识符
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// 解析字符串字面量
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// 解析整型字面量
func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
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

// 解析布尔值
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := "no prefix parse function for %s found"
	p.errors = append(p.errors, fmt.Sprintf(msg, t))
}

// 解析前缀表达式
func (p *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))
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

// 解析中缀表达式
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		// 左侧表达式
		Left: left,
	}
	// 获取当前运算符的优先级
	precedence := p.curPrecedence()
	// 读取下一个token
	p.nextToken()
	// 解析右侧表达式
	expression.Right = p.parseExpression(precedence)
	return expression
}
