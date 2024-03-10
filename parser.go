package main

import (
	"strconv"
	"yap/bytecode"
	"yap/tokens"
)

func parse(tokenList tokens.TokenList) *Env {
	env := NewEnv()

	for !tokenList.Empty() {
		parseCall(&env, &tokenList)
	}

	return &env
}

func parseCall(env *Env, tokenList *tokens.TokenList) {
	tokenList.Expect(tokens.LeftParen)

	callName := tokenList.Expect(tokens.Symbol).Content
	var argc int64 = 0

	for nextToken := tokenList.Peek(); nextToken.Type != tokens.RightParen; nextToken = tokenList.Peek() {
		argc++

		switch nextToken.Type {
		case tokens.Integer:
			t := tokenList.Consume()
			i, _ := strconv.ParseInt(t.Content, 10, 64)
			env.Push(i)
			env.Push(bytecode.Integer)
		case tokens.Float:
			t := tokenList.Consume()
			f, _ := strconv.ParseFloat(t.Content, 64)
			env.PushFloat(f)
			env.Push(bytecode.Float)
		case tokens.Symbol:
			symbolIndex := env.Symbols[nextToken.Content]
			env.Push(symbolIndex)
			env.Push(bytecode.Var)
		case tokens.LeftParen:
			parseCall(env, tokenList)
		}
	}

	env.Push(argc)

	switch callName {
	case "+":
		env.Push(bytecode.Add)
	case "yap":
		env.Push(bytecode.Print)
	default:
		fnIndex := env.Functions[callName]
		env.Push(fnIndex)
	}

	tokenList.Expect(tokens.RightParen)
}
