package main

import "github.com/ethandenny/yap/tokens"

type Scanner struct {
	tokens     tokens.TokenList
	LineNumber int
}

func scan(line string) tokens.TokenList {
	s := Scanner{
		LineNumber: 1,
	}

	for i := 0; i < len(line); i++ {
		var t tokens.Token

		switch line[i] {
		case '(':
			t = tokens.Token{
				Type:       tokens.LeftParen,
				Content:    line[i : i+1],
				LineNumber: s.LineNumber,
			}
		case ')':
			t = tokens.Token{
				Type:       tokens.RightParen,
				Content:    line[i : i+1],
				LineNumber: s.LineNumber,
			}
		case ' ':
			continue
		case '\n':
			s.LineNumber += 1
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t = scanNumber(&i, line, s.LineNumber)
		default:
			t = scanSymbol(&i, line, s.LineNumber)
		}

		s.tokens.Insert(t)
	}

	return s.tokens
}

func scanNumber(i *int, line string, LineNumber int) tokens.Token {
	start := *i
	Type := tokens.Integer

Loop:
	for ; *i < len(line); *i++ {
		switch line[*i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			if Type == tokens.Integer {
				Type = tokens.Float
			} else {
				panic("Floats aren't version numbers!")
			}
		default:
			*i--
			break Loop
		}
	}

	end := min(*i+1, len(line))

	return tokens.Token{
		Type:       Type,
		Content:    line[start:end],
		LineNumber: LineNumber,
	}
}

func scanSymbol(i *int, line string, LineNumber int) tokens.Token {
	start := *i

Loop:
	for ; *i < len(line); *i++ {
		switch line[*i] {
		case ')', '(', ' ', '\n':
			*i--
			break Loop
		}
	}

	end := min(*i+1, len(line))

	return tokens.Token{
		Type:       tokens.Symbol,
		Content:    line[start:end],
		LineNumber: LineNumber,
	}
}
