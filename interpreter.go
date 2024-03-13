package main

import (
	"fmt"

	"github.com/ethandenny/yap/tokens"
)

func eval(env *Env, stack *Stack) (int64, YapType) {
	switch popStack(stack) {
	case InstrInteger:
		return popStack(stack), IntegerT
	case InstrFloat:
		return popStack(stack), FloatT
	case InstrAdd:
		assertArgc(popStack(stack), 2)

		a, aT := eval(env, stack)
		b, bT := eval(env, stack)

		if aT == IntegerT && bT == IntegerT {
			r := a + b
			return r, IntegerT
		} else {
			var aV float64
			if aT == FloatT {
				aV = env.GetFloat(a)
			} else {
				aV = float64(a)
			}

			var bV float64
			if bT == FloatT {
				bV = env.GetFloat(b)
			} else {
				bV = float64(b)
			}

			f := aV + bV
			index := env.InsertFloat(f)

			return index, FloatT
		}
	case InstrPrint:
		argc := popStack(stack)

		if argc > 0 {
			var i int64 = 0
			for ; i < argc; i++ {
				v, t := eval(env, stack)

				if t == IntegerT {
					fmt.Print(v, " ")
				} else if t == FloatT {
					fmt.Print(env.GetFloat(v), " ")
				}
			}
		}

		fmt.Println()
	case InstrVar:
		id := popStack(stack)
		return env.GetVariable(id)
	default:
		panic("Unrecognized bytecode instruction")
	}

	return 0, NoneT
}

func evalTokens(env *Env, list *tokens.TokenList) (int64, YapType) {
	stack := parseArg(env, list)
	flipStack(&stack)
	return eval(env, &stack)
}

func assertArgc(argc int64, n int64) {
	if argc != n {
		panic("Expected different number of args")
	}
}
