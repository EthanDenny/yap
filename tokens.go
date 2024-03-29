package main

type TokenType int8

const (
	TokenNone TokenType = iota
	TokenLeftParen
	TokenRightParen
	TokenSymbol
	TokenInteger
	TokenFloat
	TokenString
	TokenLeftBracket
	TokenRightBracket
	TokenLeftBrace
	TokenRightBrace
)

type Token struct {
	Type       TokenType
	Content    string
	LineNumber int
}

type TokenList struct {
	curr   int
	tokens []Token
}

func (l *TokenList) Expect(Type TokenType) Token {
	if t := l.Consume(); t.Type == Type {
		return t
	} else {
		panic("Expected another token")
	}
}

func (l *TokenList) Insert(t Token) {
	l.tokens = append(l.tokens, t)
}

func (l *TokenList) Peek() Token {
	return l.tokens[l.curr]
}

func (l *TokenList) Consume() Token {
	l.curr++
	return l.tokens[l.curr-1]
}

func (l *TokenList) GetAll() []Token {
	return l.tokens
}

func (l *TokenList) Len() int {
	return len(l.tokens)
}

func (l *TokenList) HasToken() bool {
	return l.curr < l.Len()
}

func (l *TokenList) Append(otherList TokenList) {
	l.tokens = append(l.tokens, otherList.GetAll()...)
}
