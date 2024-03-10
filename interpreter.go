package main

import (
	"fmt"

	"github.com/ethandenny/yap/bytecode"
)

type ValueType int64

const (
	IntegerT ValueType = iota
	FloatT
	NoneT
)

func popArg(env *Env) (int64, ValueType) {
	t := eval(env)
	v := env.Pop()
	return v, t
}

func assertArgc(argc int64, n int64) {
	if argc != n {
		panic("Expected different number of args")
	}
}

func eval(env *Env) ValueType {
	switch env.Pop() {
	case bytecode.Integer:
		return IntegerT
	case bytecode.Float:
		return FloatT
	case bytecode.Add:
		assertArgc(env.Pop(), 2)

		a, aT := popArg(env)
		b, bT := popArg(env)

		if aT == IntegerT && bT == IntegerT {
			r := a + b
			env.Push(r)
			return IntegerT
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
			env.PushFloat(f)

			return FloatT
		}
	case bytecode.Print:
		argc := env.Pop()

		defer fmt.Println()

		if argc > 0 {
			argv := make([]any, argc)

			var i int64 = 0
			for ; i < argc; i++ {
				v, t := popArg(env)

				if t == IntegerT {
					argv[i] = v
				} else if t == FloatT {
					argv[i] = env.GetFloat(v)
				}
			}

			for i := len(argv) - 1; i >= 0; i-- {
				fmt.Print(argv[i], " ")
			}
		}
	default:
		panic("Unrecognized bytecode instruction")
	}

	return NoneT
}
