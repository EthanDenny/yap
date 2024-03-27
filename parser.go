package main

import (
	"strconv"
)

func parse(env *Env, symbols *SymbolTable, stack *Stack, list TokenList) {
	var tempStack Stack

	for list.HasToken() {
		switch list.Peek().Type {
		case TokenLeftParen:
			tempStack = append(tempStack, parseCall(env, symbols, &list, nil)...)
		case TokenLeftBracket:
			tempStack = append(tempStack, parseList(env, symbols, &list, nil)...)
		case TokenLeftBrace:
			tempStack = append(tempStack, parseBody(env, symbols, &list)...)
		}
	}

	flipStack(&tempStack)
	*stack = append(*stack, tempStack...)
}

func parseCall(env *Env, symbols *SymbolTable, list *TokenList, argNames []string) Stack {
	list.Expect(TokenLeftParen)

	callName := list.Expect(TokenSymbol).Content

	var stack Stack

	switch callName {
	case "=":
		varName := list.Expect(TokenSymbol).Content
		varValue, varType := evalTokens(env, symbols, list)
		env.SetVariable(symbols, varName, varValue, varType)
	case "def":
		fnName := list.Expect(TokenSymbol).Content
		var argNames []string

		list.Expect(TokenLeftBracket)
		for list.Peek().Type != TokenRightBracket {
			argName := list.Expect(TokenSymbol).Content
			argNames = append(argNames, argName)
		}
		list.Expect(TokenRightBracket)

		id := env.CreateFn(argNames)
		env.SetVariable(symbols, fnName, id, TypeFunction)

		var body TokenList

		list.Expect(TokenLeftBrace)
		braceCount := 1
		for braceCount > 0 {
			switch list.Peek().Type {
			case TokenLeftBrace:
				braceCount++
				body.Insert(list.Consume())
			case TokenRightBrace:
				braceCount--
				if braceCount > 0 {
					body.Insert(list.Consume())
				}
			default:
				body.Insert(list.Consume())
			}
		}
		list.Expect(TokenRightBrace)

		env.SetFnBody(id, body)
	default:
		var args []Stack

		for list.Peek().Type != TokenRightParen {
			args = append(args, parseArg(env, symbols, list, argNames))
		}

		var builtIns = map[string]int64{
			"+":     InstrAdd,
			"print": InstrPrint,
			"-":     InstrSub,
			"if":    InstrIf,
			"==":    InstrEq,
			"++":    InstrPush,
			"head":  InstrHead,
			"tail":  InstrTail,
			"stoi":  InstrStoi,
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

func parseList(env *Env, symbols *SymbolTable, list *TokenList, argNames []string) Stack {
	var elements []Variable

	list.Expect(TokenLeftBracket)

	for list.Peek().Type != TokenRightBracket {
		stack := parseArg(env, symbols, list, argNames)
		flipStack(&stack)
		v, t := eval(env, symbols, &stack)
		elements = append(elements, Variable{v, t})
	}

	list.Expect(TokenRightBracket)

	var id int64 = -1 // Nil
	for i := len(elements) - 1; i >= 0; i-- {
		id = env.CreateList(id, elements[i].Value, elements[i].Type)
	}

	var stack Stack

	stack = append(stack, InstrList)
	stack = append(stack, id)

	return stack
}

func parseBody(env *Env, symbols *SymbolTable, list *TokenList) Stack {
	id := env.CreateFn(make([]string, 0))

	var body TokenList

	list.Expect(TokenLeftBrace)
	braceCount := 1
	for braceCount > 0 {
		switch list.Peek().Type {
		case TokenLeftBrace:
			braceCount++
			body.Insert(list.Consume())
		case TokenRightBrace:
			braceCount--
			if braceCount > 0 {
				body.Insert(list.Consume())
			}
		default:
			body.Insert(list.Consume())
		}
	}
	list.Expect(TokenRightBrace)

	env.SetFnBody(id, body)

	var stack Stack

	stack = append(stack, InstrFn)
	stack = append(stack, id)
	stack = append(stack, 0)

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
		case "none":
			return []int64{InstrNone, 0}
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
	case TokenLeftBracket:
		return parseList(env, symbols, list, argNames)
	case TokenLeftBrace:
		return parseBody(env, symbols, list)
	default:
		panic("Unexpected token while parsing arg")
	}
}
