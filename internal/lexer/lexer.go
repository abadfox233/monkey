package lexer

import (
	"monkey/internal/token"
	"strings"
)

type Lexer struct {
	input        string
	position     int  // 当前字符的位置
	readPosition int  // 当前读取字符的位置
	ch           byte // 当前字符

}

func New(input string) *Lexer {

	l := &Lexer{input: input}

	l.readChar()

	return l

}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 读取下一个字符
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '"':
		// 读取字符串
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '=':
		if l.peekChar() == '=' {
			// 读取下一个字符
			ch := l.ch
			l.readChar()
			// 读取下一个字符
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			// 读取下一个字符
			ch := l.ch
			l.readChar()
			// 读取下一个字符
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) { // 标识符
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) { // 数字
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else { // 未知字符
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok

}

func (l *Lexer) readString() string {

	buffer := strings.Builder{}

	for {
		l.readChar()
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case '0':
				buffer.WriteString(string('\000'))
			case 'a':
				buffer.WriteString(string('\a'))
			case 'b':
				buffer.WriteString(string('\b'))
			case 'f':
				buffer.WriteString(string('\f'))
			case 'v':
				buffer.WriteString(string('\v'))
			case '"':
				buffer.WriteString(string('"'))
			case 'n':
				buffer.WriteString(string('\n'))
			case 't':
				buffer.WriteString(string('\t'))
			case 'r':
				buffer.WriteString(string('\r'))
			case '\\':
				buffer.WriteString(string('\\'))
			default:
				buffer.WriteString(string('\\'))
				buffer.WriteString(string(l.ch))
			}
		} else if l.ch == '"' || l.ch == 0 {
			break
		} else {
			buffer.WriteString(string(l.ch))
		}
	}
	return buffer.String()
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]

}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// 生成token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' { // \r 回车符 \n 换行符
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}

}
