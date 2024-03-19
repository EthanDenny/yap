package main

import (
	"fmt"
)

type Arg struct {
	v int64
	t YapType
}

func eval(env *Env, stack *Stack, args []Arg) (int64, YapType) {
	switch popStack(stack) {
	case InstrInteger:
		return popStack(stack), TypeInteger
	case InstrFloat:
		return popStack(stack), TypeFloat
	case InstrString:
		return popStack(stack), TypeString
	case InstrAdd:
		assertArgc(popStack(stack), 2)

		a, aT := eval(env, stack, args)
		b, bT := eval(env, stack, args)

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

		a, aT := eval(env, stack, args)
		b, bT := eval(env, stack, args)

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

		a, aT := eval(env, stack, args)
		b, bT := eval(env, stack, args)

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

		return result, TypeBool
	case InstrPrint:
		argc := popStack(stack)

		var i int64 = 0
		for ; i < argc; i++ {
			v, t := eval(env, stack, args)

			switch t {
			case TypeInteger:
				fmt.Print(v, " ")
			case TypeFloat:
				fmt.Print(env.GetFloat(v), " ")
			case TypeString:
				fmt.Print(env.GetString(v), " ")
			case TypeBool:
				if v == 1 {
					fmt.Print("true")
				} else {
					fmt.Print("false")
				}
			}
		}

		fmt.Println()
	case InstrVar:
		id := popStack(stack)
		return env.GetVariable(id)
	case InstrFn:
		id := popStack(stack)
		argc, body := env.GetFn(id)
		assertArgc(popStack(stack), argc)

		var fnArgs []Arg

		var i int64 = 0
		for ; i < argc; i++ {
			v, t := eval(env, stack, args)
			fnArgs = append(fnArgs, Arg{
				v: v,
				t: t,
			})
		}

		var returnV int64
		var returnT YapType

		for len(body) > 0 {
			returnV, returnT = eval(env, &body, fnArgs)
		}

		return returnV, returnT
	case InstrArg:
		argN := popStack(stack)
		arg := args[argN]
		return arg.v, arg.t
	case InstrIf:
		assertArgc(popStack(stack), 3)

		pred, predT := eval(env, stack, args)

		if predT != TypeBool {
			panic("Need boolean for predicate")
		}

		if pred == 1 {
			v, vT := eval(env, stack, args)
			popArg(env, stack)
			return v, vT
		} else {
			popArg(env, stack)
			v, vT := eval(env, stack, args)
			return v, vT
		}
	case InstrBool:
		return popStack(stack), TypeBool
	case InstrPush:
		assertArgc(popStack(stack), 2)

		e, eT := eval(env, stack, args)
		l, lT := eval(env, stack, args)

		if eT == TypeString && lT == TypeString {
			eS := env.GetString(e)
			lS := env.GetString(l)
			index := env.InsertString(eS + lS)
			return index, TypeString
		}
	case InstrHead:
		assertArgc(popStack(stack), 1)

		l, lT := eval(env, stack, args)

		if lT == TypeString {
			lS := env.GetString(l)
			var head string
			if len(lS) > 1 {
				head = lS[0:1]
			} else {
				head = lS
			}
			index := env.InsertString(head)
			return index, TypeString
		}
	case InstrTail:
		assertArgc(popStack(stack), 1)

		l, lT := eval(env, stack, args)

		if lT == TypeString {
			lS := env.GetString(l)
			var tail string
			if len(lS) > 1 {
				tail = lS[1:]
			} else {
				tail = ""
			}
			index := env.InsertString(tail)
			return index, TypeString
		}
	default:
		panic("Unrecognized instruction")
	}

	return 0, TypeNone
}

func popArg(env *Env, stack *Stack) {
	switch popStack(stack) {
	case InstrInteger, InstrFloat, InstrString, InstrVar, InstrArg, InstrBool:
		popStack(stack)
	case InstrFn:
		id := popStack(stack)
		argc, _ := env.GetFn(id)
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

func evalTokens(env *Env, list *TokenList) (int64, YapType) {
	stack := parseArg(env, list, nil)
	flipStack(&stack)
	return eval(env, &stack, nil)
}

func assertArgc(argc int64, n int64) {
	if argc != n {
		panic("Expected different number of args")
	}
}

func disassemble(env *Env, stack *Stack, args []Arg) {
	switch popStack(stack) {
	case InstrInteger:
		fmt.Print(popStack(stack))
	case InstrFloat:
		fmt.Print(env.GetFloat(popStack(stack)))
	case InstrString:
		fmt.Print("\"", env.GetString(popStack(stack)), "\"")
	case InstrAdd:
		popStack(stack)
		fmt.Print("(ADD ")
		disassemble(env, stack, args)
		fmt.Print(" ")
		disassemble(env, stack, args)
		fmt.Print(")")
	case InstrSub:
		popStack(stack)
		fmt.Print("(SUB ")
		disassemble(env, stack, args)
		fmt.Print(" ")
		disassemble(env, stack, args)
		fmt.Print(")")
	case InstrPush:
		popStack(stack)
		fmt.Print("(PUSH ")
		disassemble(env, stack, args)
		fmt.Print(" ")
		disassemble(env, stack, args)
		fmt.Print(")")
	case InstrEq:
		popStack(stack)
		fmt.Print("(EQ ")
		disassemble(env, stack, args)
		fmt.Print(" ")
		disassemble(env, stack, args)
		fmt.Print(")")
	case InstrPrint:
		argc := popStack(stack)
		fmt.Print("(PRINT ")
		var i int64 = 0
		for ; i < argc; i++ {
			disassemble(env, stack, args)
			if i+1 < argc {
				fmt.Print(" ")
			}
		}
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
		var i int64 = 0
		for ; i < argc; i++ {
			disassemble(env, stack, args)
			if i+1 < argc {
				fmt.Print(" ")
			}
		}
		fmt.Print(")")
	case InstrArg:
		argN := popStack(stack)
		fmt.Print("(ARG ", argN, " => ")

		arg := args[argN]
		switch arg.t {
		case TypeInteger:
			fmt.Print(arg.v)
		case TypeFloat:
			fmt.Print(env.GetFloat(arg.v))
		case TypeString:
			fmt.Print("\"", env.GetString(arg.v), "\"")
		case TypeBool:
			if arg.v == 1 {
				fmt.Print("true")
			} else {
				fmt.Print("false")
			}
		}

		fmt.Print(")")
	case InstrIf:
		popStack(stack)
		fmt.Print("(IF ")
		disassemble(env, stack, args)
		fmt.Print(" THEN ")
		disassemble(env, stack, args)
		fmt.Print(" ELSE ")
		disassemble(env, stack, args)
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
		disassemble(env, stack, args)
		fmt.Print(")")
	case InstrTail:
		popStack(stack)
		fmt.Print("(TAIL ")
		disassemble(env, stack, args)
		fmt.Print(")")
	default:
	}
}
