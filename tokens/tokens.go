package tokens

type TokenType int8

const (
	None TokenType = iota
	LeftParen
	RightParen
	Symbol
	Integer
	Float
)

type Token struct {
	Type       TokenType
	Content    string
	LineNumber int
}

func none() Token {
	return Token{
		Type:       None,
		Content:    "()",
		LineNumber: 0,
	}
}

type TokenList struct {
	curr   int
	tokens []Token
}

func (list *TokenList) Expect(Type TokenType) Token {
	if t := list.Consume(); t.Type == Type {
		return t
	} else {
		panic("Expected another token")
	}
}

func (list *TokenList) Insert(t Token) {
	list.tokens = append(list.tokens, t)
}

func (list TokenList) Peek() Token {
	return list.tokens[list.curr]
}

func (list *TokenList) Consume() Token {
	list.curr += 1
	return list.tokens[list.curr-1]
}

func (list TokenList) GetAll() []Token {
	return list.tokens
}

func (list TokenList) Len() int {
	return len(list.tokens)
}

func (list TokenList) Empty() bool {
	return list.Len() == list.curr
}
