package main

import (
	"strconv"

	"github.com/ethandenny/yap/tokens"
)

func parse(env *Env, stack *Stack, list tokens.TokenList) {
	var tempStack Stack
	for list.HasToken() {
		tempStack = append(tempStack, parseCall(env, &list, nil)...)
	}
	flipStack(&tempStack)

	*stack = append(*stack, tempStack...)
}

func parseCall(env *Env, list *tokens.TokenList, argNames []string) Stack {
	list.Expect(tokens.LeftParen)

	callName := list.Expect(tokens.Symbol).Content

	var stack Stack

	switch callName {
	case "let":
		varName := list.Expect(tokens.Symbol).Content
		varValue, varType := evalTokens(env, list)

		env.SetVariable(varName, varValue, varType)
	case "def":
		fnName := list.Expect(tokens.Symbol).Content
		var argNames []string

		list.Expect(tokens.LeftParen)
		for list.Peek().Type != tokens.RightParen {
			argName := list.Expect(tokens.Symbol).Content
			argNames = append(argNames, argName)
		}
		list.Expect(tokens.RightParen)

		var fnBody Stack

		list.Expect(tokens.LeftParen)
		for list.Peek().Type != tokens.RightParen {
			nextCall := parseArg(env, list, argNames)
			fnBody = append(fnBody, nextCall...)
		}
		list.Expect(tokens.RightParen)

		id := env.CreateFn(int64(len(argNames)), fnBody)
		env.SetVariable(fnName, id, FunctionT)
	default:
		var args []Stack

		for list.Peek().Type != tokens.RightParen {
			args = append(args, parseArg(env, list, argNames))
		}

		var builtIns = map[string]int64{
			"+":   InstrAdd,
			"yap": InstrPrint,
		}

		if f, containsKey := builtIns[callName]; containsKey {
			stack = append(stack, f)
		} else {
			id := env.GetSymbol(callName)

			stack = append(stack, InstrFn)
			stack = append(stack, id)
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

func parseArg(env *Env, list *tokens.TokenList, argNames []string) Stack {
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

		for id, name := range argNames {
			if name == t.Content {
				return []int64{InstrArg, int64(id)}
			}
		}

		id := env.GetSymbol(t.Content)
		return []int64{InstrVar, id}
	case tokens.LeftParen:
		return parseCall(env, list, argNames)
	default:
		panic("Unexpected token while parsing arg")
	}
}
