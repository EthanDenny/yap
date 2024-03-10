package main

import "fmt"

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
	case Integer:
		return IntegerT
	case Float:
		return FloatT
	case Add:
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
			index := env.InsertFloat(f)
			env.Push(index)

			return FloatT
		}
	case Print:
		argc := env.Pop()

		if argc > 0 {
			var i int64 = 0
			for ; i < argc; i++ {
				v, t := popArg(env)

				if t == IntegerT {
					fmt.Print(v, " ")
				} else if t == FloatT {
					fmt.Print(env.GetFloat(v), " ")
				}
			}
		}

		fmt.Println()
	default:
		panic("Unrecognized bytecode instruction")
	}

	return NoneT
}
