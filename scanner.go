package main

type Scanner struct {
	innerStr   string
	strPos     int
	lineNumber int
	tokens     TokenList
}

func (s *Scanner) getChar() byte {
	return s.innerStr[s.strPos]
}

func (s *Scanner) getNextChar() byte {
	return s.innerStr[s.strPos+1]
}

func (s *Scanner) Len() int {
	return len(s.innerStr) - s.strPos
}

func (s *Scanner) scanSingleChar(Type TokenType) {
	s.tokens.Insert(Token{
		Type:       Type,
		Content:    s.innerStr[s.strPos : s.strPos+1],
		LineNumber: s.lineNumber,
	})
}

func (s *Scanner) scanNumber() {
	start := s.strPos
	Type := TokenInteger

	if s.getChar() == '-' {
		s.strPos++
	}

Loop:
	for ; s.strPos < len(s.innerStr); s.strPos++ {
		switch s.getChar() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			if Type == TokenInteger {
				Type = TokenFloat
			} else {
				panic("Floats aren't version numbers!")
			}
		default:
			s.strPos--
			break Loop
		}
	}

	end := min(s.strPos+1, len(s.innerStr))

	s.tokens.Insert(Token{
		Type:       Type,
		Content:    s.innerStr[start:end],
		LineNumber: s.lineNumber,
	})
}

func (s *Scanner) scanString() {
	s.strPos++
	start := s.strPos
	foundClosingQuote := false

Loop:
	for ; s.strPos < len(s.innerStr); s.strPos++ {
		switch s.getChar() {
		case '"':
			foundClosingQuote = true
			break Loop
		}
	}

	if !foundClosingQuote {
		panic("Expected closing quote, found EOF")
	}

	s.tokens.Insert(Token{
		Type:       TokenString,
		Content:    s.innerStr[start:s.strPos],
		LineNumber: s.lineNumber,
	})
}

func (s *Scanner) scanSymbol() {
	start := s.strPos

Loop:
	for ; s.strPos < len(s.innerStr); s.strPos++ {
		switch s.getChar() {
		case '(', ')', '[', ']', ' ', '\n':
			s.strPos--
			break Loop
		}
	}

	end := min(s.strPos+1, len(s.innerStr))

	s.tokens.Insert(Token{
		Type:       TokenSymbol,
		Content:    s.innerStr[start:end],
		LineNumber: s.lineNumber,
	})
}

func scan(str string) TokenList {
	s := Scanner{
		innerStr:   str,
		lineNumber: 1,
	}

	for ; s.strPos < len(str); s.strPos++ {
		switch s.getChar() {
		case '(':
			s.scanSingleChar(TokenLeftParen)
		case ')':
			s.scanSingleChar(TokenRightParen)
		case ' ':
			continue
		case '\n':
			s.lineNumber++
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s.scanNumber()
		case '-':
			if s.Len() > 1 {
				switch s.getNextChar() {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					s.scanNumber()
				default:
					s.scanSymbol()
				}
			} else {
				s.scanSymbol()
			}
		case '"':
			s.scanString()
		default:
			s.scanSymbol()
		}
	}

	return s.tokens
}
