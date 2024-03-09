package main

import (
	"fmt"
	"yap/bytecode"
)

type ValueType int64

const (
	IntegerT ValueType = iota
	FloatT
	StringT
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

		a, a_t := popArg(env)
		b, b_t := popArg(env)

		if a_t == IntegerT && b_t == IntegerT {
			r := a + b
			env.Push(r)
			return IntegerT
		} else {
			var a_v float64
			if a_t == FloatT {
				a_v = env.GetFloat(a)
			} else {
				a_v = float64(a)
			}

			var b_v float64
			if b_t == FloatT {
				b_v = env.GetFloat(b)
			} else {
				b_v = float64(b)
			}

			f := a_v + b_v
			env.PushFloat(f)

			return FloatT
		}
	case bytecode.Print:
		argc := env.Pop()

		defer fmt.Println()

		if argc > 0 {
			var i int64 = 0
			for ; i < argc; i++ {
				v, t := popArg(env)
				var v_out interface{}

				if t == IntegerT {
					v_out = v
				} else if t == FloatT {
					v_out = env.GetFloat(v)
				}

				defer fmt.Print(v_out, " ")
			}
		}
	}

	return NoneT
}
