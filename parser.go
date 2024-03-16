package main

import (
	"strconv"
)

func parse(env *Env, stack *Stack, list TokenList) {
	var tempStack Stack
	for list.HasToken() {
		tempStack = append(tempStack, parseCall(env, &list, nil)...)
	}
	flipStack(&tempStack)

	*stack = append(*stack, tempStack...)
}

func parseCall(env *Env, list *TokenList, argNames []string) Stack {
	list.Expect(TokenLeftParen)

	callName := list.Expect(TokenSymbol).Content

	var stack Stack

	switch callName {
	case "let":
		varName := list.Expect(TokenSymbol).Content
		varValue, varType := evalTokens(env, list)

		env.SetVariable(varName, varValue, varType)
	case "def":
		fnName := list.Expect(TokenSymbol).Content
		var argNames []string

		list.Expect(TokenLeftParen)
		for list.Peek().Type != TokenRightParen {
			argName := list.Expect(TokenSymbol).Content
			argNames = append(argNames, argName)
		}
		list.Expect(TokenRightParen)

		id := env.CreateFn(int64(len(argNames)))
		env.SetVariable(fnName, id, TypeFunction)

		var fnBody Stack

		list.Expect(TokenLeftParen)
		for list.Peek().Type != TokenRightParen {
			nextCall := parseArg(env, list, argNames)
			fnBody = append(fnBody, nextCall...)
		}
		list.Expect(TokenRightParen)

		env.SetFnBody(id, fnBody)
	default:
		var args []Stack

		for list.Peek().Type != TokenRightParen {
			args = append(args, parseArg(env, list, argNames))
		}

		var builtIns = map[string]int64{
			"+":   InstrAdd,
			"yap": InstrPrint,
			"-":   InstrSub,
			"if":  InstrIf,
			"=":   InstrEq,
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

	list.Expect(TokenRightParen)

	return stack
}

func parseArg(env *Env, list *TokenList, argNames []string) Stack {
	nextToken := list.Peek()

	switch nextToken.Type {
	case TokenInteger:
		t := list.Consume()
		i, _ := strconv.ParseInt(t.Content, 10, 64)
		return []int64{InstrInteger, i}
	case TokenFloat:
		t := list.Consume()
		f, _ := strconv.ParseFloat(t.Content, 64)
		index := env.InsertFloat(f)
		return []int64{InstrFloat, index}
	case TokenSymbol:
		t := list.Consume()

		for id, name := range argNames {
			if name == t.Content {
				return []int64{InstrArg, int64(id)}
			}
		}

		switch t.Content {
		case "true":
			return []int64{InstrBool, 1}
		case "false":
			return []int64{InstrBool, 0}
		default:
			id := env.GetSymbol(t.Content)
			return []int64{InstrVar, id}
		}
	case TokenLeftParen:
		return parseCall(env, list, argNames)
	default:
		panic("Unexpected token while parsing arg")
	}
}
