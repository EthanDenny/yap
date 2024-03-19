package main

import (
	"strconv"
)

func parse(env *Env, symbols *SymbolTable, stack *Stack, list TokenList) {
	var tempStack Stack
	for list.HasToken() {
		tempStack = append(tempStack, parseCall(env, symbols, &list, nil)...)
	}
	flipStack(&tempStack)

	*stack = append(*stack, tempStack...)
}

func parseCall(env *Env, symbols *SymbolTable, list *TokenList, argNames []string) Stack {
	list.Expect(TokenLeftParen)

	callName := list.Expect(TokenSymbol).Content

	var stack Stack

	switch callName {
	case "let":
		varName := list.Expect(TokenSymbol).Content
		varValue, varType := evalTokens(env, symbols, list)
		env.SetVariable(symbols, varName, varValue, varType)
	case "def":
		fnName := list.Expect(TokenSymbol).Content
		var argNames []string

		list.Expect(TokenLeftParen)
		for list.Peek().Type != TokenRightParen {
			argName := list.Expect(TokenSymbol).Content
			argNames = append(argNames, argName)
		}
		list.Expect(TokenRightParen)

		id := env.CreateFn(argNames)
		env.SetVariable(symbols, fnName, id, TypeFunction)

		var fnBody TokenList

		list.Expect(TokenLeftParen)
		parenCount := 1
		for parenCount > 0 {
			switch list.Peek().Type {
			case TokenLeftParen:
				parenCount++
				fnBody.Insert(list.Consume())
			case TokenRightParen:
				parenCount--
				if parenCount > 0 {
					fnBody.Insert(list.Consume())
				}
			default:
				fnBody.Insert(list.Consume())
			}
		}
		list.Expect(TokenRightParen)

		env.SetFnBody(id, fnBody)
	default:
		var args []Stack

		for list.Peek().Type != TokenRightParen {
			args = append(args, parseArg(env, symbols, list, argNames))
		}

		var builtIns = map[string]int64{
			"+":    InstrAdd,
			"yap":  InstrPrint,
			"-":    InstrSub,
			"if":   InstrIf,
			"=":    InstrEq,
			"push": InstrPush,
			"head": InstrHead,
			"tail": InstrTail,
		}

		if f, containsKey := builtIns[callName]; containsKey {
			stack = append(stack, f)
		} else {
			id := symbols.Get(callName)

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

func parseArg(env *Env, symbols *SymbolTable, list *TokenList, argNames []string) Stack {
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
	case TokenString:
		t := list.Consume()
		index := env.InsertString(t.Content)
		return []int64{InstrString, index}
	case TokenSymbol:
		t := list.Consume()

		for id, name := range argNames {
			if name == t.Content {
				return []int64{InstrVar, int64(id)}
			}
		}

		switch t.Content {
		case "true":
			return []int64{InstrBool, 1}
		case "false":
			return []int64{InstrBool, 0}
		default:
			id := symbols.Get(t.Content)
			return []int64{InstrVar, id}
		}
	case TokenLeftParen:
		return parseCall(env, symbols, list, argNames)
	default:
		panic("Unexpected token while parsing arg")
	}
}
