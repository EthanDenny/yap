package main

import (
	"strconv"

	"github.com/ethandenny/yap/tokens"
)

func parse(env *Env, stack *Stack, list tokens.TokenList) {
	var tempStack Stack
	for list.HasToken() {
		tempStack = append(tempStack, parseCall(env, &list)...)
	}
	flipStack(&tempStack)

	*stack = append(*stack, tempStack...)
}

func parseCall(env *Env, list *tokens.TokenList) Stack {
	list.Expect(tokens.LeftParen)

	callName := list.Expect(tokens.Symbol).Content

	var stack Stack

	switch callName {
	case "let":
		t := list.Consume()
		if t.Type != tokens.Symbol {
			panic("Expected symbol")
		}
		varName := t.Content
		varValue, varType := evalTokens(env, list)

		env.SetVariable(varName, varValue, varType)
	default:
		var args []Stack

		for list.Peek().Type != tokens.RightParen {
			args = append(args, parseArg(env, list))
		}

		var functions = map[string]int64{
			"+":   InstrAdd,
			"yap": InstrPrint,
		}

		if f, containsKey := functions[callName]; containsKey {
			stack = append(stack, f)
		} else {
			fnIndex := env.GetFn(callName)
			stack = append(stack, fnIndex)
		}

		stack = append(stack, int64(len(args)))

		for _, arg := range args {
			for _, instr := range arg {
				stack = append(stack, instr)
			}
		}
	}

	list.Expect(tokens.RightParen)

	return stack
}

func parseArg(env *Env, list *tokens.TokenList) Stack {
	nextToken := list.Peek()

	switch nextToken.Type {
	case tokens.Integer:
		t := list.Consume()
		i, _ := strconv.ParseInt(t.Content, 10, 64)
		return []int64{InstrInteger, i}
	case tokens.Float:
		t := list.Consume()
		f, _ := strconv.ParseFloat(t.Content, 64)
		index := env.InsertFloat(f)
		return []int64{InstrFloat, index}
	case tokens.Symbol:
		t := list.Consume()
		id := env.GetSymbol(t.Content)
		return []int64{InstrVar, id}
	case tokens.LeftParen:
		return parseCall(env, list)
	default:
		panic("Unexpected token while parsing arg")
	}
}
