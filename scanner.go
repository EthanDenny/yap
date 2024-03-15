package main

import "github.com/ethandenny/yap/tokens"

type Scanner struct {
	innerStr   string
	strPos     int
	lineNumber int
	tokens     tokens.TokenList
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

func (s *Scanner) scanSingleChar(Type tokens.TokenType) {
	s.tokens.Insert(tokens.Token{
		Type:       Type,
		Content:    s.innerStr[s.strPos : s.strPos+1],
		LineNumber: s.lineNumber,
	})
}

func (s *Scanner) scanNumber() {
	start := s.strPos
	Type := tokens.Integer

	if s.getChar() == '-' {
		s.strPos++
	}

Loop:
	for ; s.strPos < len(s.innerStr); s.strPos++ {
		switch s.getChar() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			if Type == tokens.Integer {
				Type = tokens.Float
			} else {
				panic("Floats aren't version numbers!")
			}
		default:
			s.strPos--
			break Loop
		}
	}

	end := min(s.strPos+1, len(s.innerStr))

	s.tokens.Insert(tokens.Token{
		Type:       Type,
		Content:    s.innerStr[start:end],
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

	s.tokens.Insert(tokens.Token{
		Type:       tokens.Symbol,
		Content:    s.innerStr[start:end],
		LineNumber: s.lineNumber,
	})
}

func scan(str string) tokens.TokenList {
	s := Scanner{
		innerStr:   str,
		lineNumber: 1,
	}

	for ; s.strPos < len(str); s.strPos++ {
		switch s.getChar() {
		case '(':
			s.scanSingleChar(tokens.LeftParen)
		case ')':
			s.scanSingleChar(tokens.RightParen)
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
		default:
			s.scanSymbol()
		}
	}

	return s.tokens
}
