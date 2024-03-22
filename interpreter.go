package main

import (
	"fmt"
)

func eval(env *Env, symbols *SymbolTable, stack *Stack) (int64, YapType) {
	switch popStack(stack) {
	case InstrInteger:
		return popStack(stack), TypeInteger
	case InstrFloat:
		return popStack(stack), TypeFloat
	case InstrString:
		return popStack(stack), TypeString
	case InstrAdd:
		assertArgc(popStack(stack), 2)

		a, aT := eval(env, symbols, stack)
		b, bT := eval(env, symbols, stack)

		if aT == TypeInteger && bT == TypeInteger {
			r := a + b
			return r, TypeInteger
		} else {
			var aV float64
			if aT == TypeFloat {
				aV = env.GetFloat(a)
			} else {
				aV = float64(a)
			}

			var bV float64
			if bT == TypeFloat {
				bV = env.GetFloat(b)
			} else {
				bV = float64(b)
			}

			f := aV + bV
			index := env.InsertFloat(f)

			return index, TypeFloat
		}
	case InstrSub:
		assertArgc(popStack(stack), 2)

		a, aT := eval(env, symbols, stack)
		b, bT := eval(env, symbols, stack)

		if aT == TypeInteger && bT == TypeInteger {
			r := a - b
			return r, TypeInteger
		} else {
			var aV float64
			if aT == TypeFloat {
				aV = env.GetFloat(a)
			} else {
				aV = float64(a)
			}

			var bV float64
			if bT == TypeFloat {
				bV = env.GetFloat(b)
			} else {
				bV = float64(b)
			}

			f := aV - bV
			index := env.InsertFloat(f)

			return index, TypeFloat
		}
	case InstrEq:
		assertArgc(popStack(stack), 2)

		a, aT := eval(env, symbols, stack)
		b, bT := eval(env, symbols, stack)

		var result int64 = 0

		if aT == TypeInteger && bT == TypeInteger {
			if a == b {
				result = 1
			}
		}
		if aT == TypeInteger && bT == TypeFloat {
			f := env.GetFloat(b)
			if a == int64(f) {
				result = 1
			}
		}
		if aT == TypeFloat && bT == TypeInteger {
			f := env.GetFloat(a)
			if int64(f) == b {
				result = 1
			}
		}
		if aT == TypeFloat && bT == TypeFloat {
			aF := env.GetFloat(a)
			bF := env.GetFloat(b)
			if aF == bF {
				result = 1
			}
		}
		if aT == TypeString && bT == TypeString {
			aS := env.GetString(a)
			bS := env.GetString(b)
			if aS == bS {
				result = 1
			}
		}
		if aT == TypeNone && bT == TypeNone {
			result = 1
		}

		return result, TypeBool
	case InstrPrint:
		argc := popStack(stack)

		var i int64 = 0
		for ; i < argc; i++ {
			v, t := eval(env, symbols, stack)
			printArg(env, v, t)
			fmt.Print(" ")
		}

		fmt.Println()
	case InstrVar:
		id := popStack(stack)
		return env.GetVariable(id)
	case InstrFn:
		id := popStack(stack)
		argNames, body := env.GetFn(id)
		argc := int64(len(argNames))
		assertArgc(popStack(stack), argc)

		fnSymbols := NewSymbolTable(symbols)

		var i int64 = 0
		for ; i < argc; i++ {
			v, t := eval(env, symbols, stack)
			env.SetVariable(fnSymbols, argNames[i], v, t)

		}

		var returnV int64
		var returnT YapType

		var bodyStack Stack
		parse(env, fnSymbols, &bodyStack, body)

		for len(bodyStack) > 0 {
			returnV, returnT = eval(env, fnSymbols, &bodyStack)
		}

		return returnV, returnT
	case InstrIf:
		assertArgc(popStack(stack), 3)

		pred, predT := eval(env, symbols, stack)

		if predT != TypeBool {
			panic("Need boolean for predicate")
		}

		if pred == 1 {
			v, vT := eval(env, symbols, stack)
			popArg(env, stack)
			return v, vT
		} else {
			popArg(env, stack)
			v, vT := eval(env, symbols, stack)
			return v, vT
		}
	case InstrBool:
		return popStack(stack), TypeBool
	case InstrPush:
		assertArgc(popStack(stack), 2)

		e, eT := eval(env, symbols, stack)
		l, lT := eval(env, symbols, stack)

		if eT == TypeString && lT == TypeString {
			eS := env.GetString(e)
			lS := env.GetString(l)
			index := env.InsertString(eS + lS)
			return index, TypeString
		} else if lT == TypeList {
			id := env.CreateList(l, e, eT)
			return id, TypeList
		}
	case InstrHead:
		assertArgc(popStack(stack), 1)

		l, lT := eval(env, symbols, stack)

		switch lT {
		case TypeString:
			lS := env.GetString(l)
			var head string
			if len(lS) > 1 {
				head = lS[0:1]
			} else {
				head = lS
			}
			index := env.InsertString(head)
			return index, TypeString
		case TypeList:
			if l != -1 {
				head, _ := env.GetList(l)
				v, t := env.GetVariable(head)
				return v, t
			}
		}
	case InstrTail:
		assertArgc(popStack(stack), 1)

		l, lT := eval(env, symbols, stack)

		switch lT {
		case TypeString:
			lS := env.GetString(l)
			var tail string
			if len(lS) > 1 {
				tail = lS[1:]
			} else {
				tail = ""
			}
			index := env.InsertString(tail)
			return index, TypeString
		case TypeList:
			_, tail := env.GetList(l)
			if tail != -1 {
				return tail, TypeList
			}
		}
	case InstrList:
		argc := popStack(stack)

		var elements []Variable

		var i int64 = 0
		for ; i < argc; i++ {
			v, t := eval(env, symbols, stack)
			elements = append(elements, Variable{v, t})
		}

		var nextID int64 = -1 // Nil
		for i := len(elements) - 1; i >= 0; i-- {
			nextID = env.CreateList(nextID, elements[i].Value, elements[i].Type)
		}

		return nextID, TypeList
	case InstrNone:
		popStack(stack)
		return 0, TypeNone
	default:
		panic("Unrecognized instruction")
	}

	return 0, TypeNone
}

func popArg(env *Env, stack *Stack) {
	instr := popStack(stack)
	switch instr {
	case InstrInteger, InstrFloat, InstrString, InstrVar, InstrBool, InstrNone:
		popStack(stack)
	case InstrFn:
		id := popStack(stack)
		argNames, _ := env.GetFn(id)
		argc := int64(len(argNames))
		assertArgc(popStack(stack), argc)
		var i int64 = 0
		for ; i < argc; i++ {
			popArg(env, stack)
		}
	default:
		argc := popStack(stack)
		var i int64 = 0
		for ; i < argc; i++ {
			popArg(env, stack)
		}
	}
}

func evalTokens(env *Env, symbols *SymbolTable, list *TokenList) (int64, YapType) {
	stack := parseArg(env, symbols, list, nil)
	flipStack(&stack)
	return eval(env, symbols, &stack)
}

func assertArgc(argc int64, n int64) {
	if argc != n {
		panic("Expected different number of args")
	}
}

func disassemble(env *Env, stack *Stack) {
	stackCopy := *stack
	dis(env, &stackCopy)
	fmt.Println()
}

func dis(env *Env, stack *Stack) {
	switch popStack(stack) {
	case InstrInteger:
		fmt.Print(popStack(stack))
	case InstrFloat:
		fmt.Print(env.GetFloat(popStack(stack)))
	case InstrString:
		fmt.Print("\"", env.GetString(popStack(stack)), "\"")
	case InstrNone:
		fmt.Print("NONE")
	case InstrAdd:
		popStack(stack)
		fmt.Print("(ADD ")
		dis(env, stack)
		fmt.Print(" ")
		dis(env, stack)
		fmt.Print(")")
	case InstrSub:
		popStack(stack)
		fmt.Print("(SUB ")
		dis(env, stack)
		fmt.Print(" ")
		dis(env, stack)
		fmt.Print(")")
	case InstrPush:
		popStack(stack)
		fmt.Print("(PUSH ")
		dis(env, stack)
		fmt.Print(" ")
		dis(env, stack)
		fmt.Print(")")
	case InstrEq:
		popStack(stack)
		fmt.Print("(EQ ")
		dis(env, stack)
		fmt.Print(" ")
		dis(env, stack)
		fmt.Print(")")
	case InstrPrint:
		argc := popStack(stack)
		fmt.Print("(PRINT ")
		disArgs(env, stack, argc)
		fmt.Print(")")
	case InstrVar:
		fmt.Print("(VAR ", popStack(stack), ")")
	case InstrFn:
		id := popStack(stack)
		argc := popStack(stack)
		fmt.Print("(FN#", id)
		if argc > 0 {
			fmt.Print(" ")
		}
		disArgs(env, stack, argc)
		fmt.Print(")")
	case InstrIf:
		popStack(stack)
		fmt.Print("(IF ")
		dis(env, stack)
		fmt.Print(" THEN ")
		dis(env, stack)
		fmt.Print(" ELSE ")
		dis(env, stack)
		fmt.Print(")")
	case InstrBool:
		if popStack(stack) == 1 {
			fmt.Print("true")
		} else {
			fmt.Print("false")
		}
	case InstrHead:
		popStack(stack)
		fmt.Print("(HEAD ")
		dis(env, stack)
		fmt.Print(")")
	case InstrTail:
		popStack(stack)
		fmt.Print("(TAIL ")
		dis(env, stack)
		fmt.Print(")")
	case InstrList:
		argc := popStack(stack)
		fmt.Print("(LIST ")
		disArgs(env, stack, argc)
		fmt.Print(")")
	default:
	}
}

func disArgs(env *Env, stack *Stack, argc int64) {
	var i int64 = 0
	for ; i < argc; i++ {
		dis(env, stack)
		if i+1 < argc {
			fmt.Print(" ")
		}
	}
}

func printArg(env *Env, v int64, t YapType) {
	switch t {
	case TypeInteger:
		fmt.Print(v)
	case TypeFloat:
		fmt.Print(env.GetFloat(v))
	case TypeNone:
		fmt.Print("none")
	case TypeString:
		fmt.Print("\"", env.GetString(v), "\"")
	case TypeBool:
		if v == 1 {
			fmt.Print("true")
		} else {
			fmt.Print("false")
		}
	case TypeList:
		head, tail := env.GetList(v)
		v, t := env.GetVariable(head)
		fmt.Print("[")
		printArg(env, v, t)
		for tail != -1 {
			head, tail = env.GetList(tail)
			v, t = env.GetVariable(head)
			fmt.Print(", ")
			printArg(env, v, t)
		}
		fmt.Print("]")
	}
}
