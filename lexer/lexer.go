package lexer

import (
	"strings"

	"github.com/jcarley/bashscript/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (this *Lexer) NextToken() token.Token {
	var tok token.Token

	this.skipWhitespace()

	switch this.ch {
	case '=':
		tok = newToken(token.ASSIGN, this.ch)
	case '+':
		tok = newToken(token.PLUS, this.ch)
	case '-':
		tok = newToken(token.MINUS, this.ch)
	case '!':
		tok = newToken(token.BANG, this.ch)
	case '/':
		tok = newToken(token.SLASH, this.ch)
	case '*':
		tok = newToken(token.ASTERISK, this.ch)
	case '<':
		tok = newToken(token.LT, this.ch)
	case '>':
		tok = newToken(token.GT, this.ch)
	case ';':
		tok = newToken(token.SEMICOLON, this.ch)
	case ',':
		tok = newToken(token.COMMA, this.ch)
	case '(':
		tok = newToken(token.LPAREN, this.ch)
	case ')':
		tok = newToken(token.RPAREN, this.ch)
	case '{':
		tok = newToken(token.LBRACE, this.ch)
	case '}':
		tok = newToken(token.RBRACE, this.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(this.ch) {
			tok.Literal = this.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(this.ch) {
			tok.Literal = this.readNumber()
			tok.Type = numberToken(tok.Literal)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, this.ch)
		}
	}

	this.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (this *Lexer) readChar() {
	if this.readPosition >= len(this.input) {
		this.ch = 0
	} else {
		this.ch = this.input[this.readPosition]
	}
	this.position = this.readPosition
	this.readPosition += 1
}

func (this *Lexer) readIdentifier() string {
	position := this.position
	for isLetter(this.ch) {
		this.readChar()
	}
	return this.input[position:this.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (this *Lexer) readNumber() string {
	position := this.position
	for isDigit(this.ch) {
		this.readChar()
	}
	return this.input[position:this.position]
}

func (this *Lexer) skipWhitespace() {
	for this.ch == ' ' || this.ch == '\t' || this.ch == '\n' || this.ch == '\r' {
		this.readChar()
	}
}

func (this *Lexer) peekChar() byte {
	if this.readPosition >= len(this.input) {
		return 0
	} else {
		return this.input[this.readPosition]
	}
}

func numberToken(literal string) token.TokenType {
	if isFloat(literal) {
		return token.FLOAT
	} else {
		return token.INT
	}
	return token.ILLEGAL
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}

func isFloat(literal string) bool {
	return strings.Contains(literal, ".")
}
