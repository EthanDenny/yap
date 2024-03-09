package main

type Env struct {
	Stack []int64

	Functions map[string]int64
	fnIndex   int64

	Symbols     map[string]int64
	symbolIndex int64

	floats     map[int64]float64
	floatIndex int64
}

func NewEnv() Env {
	return Env{
		Stack: make([]int64, 0, 1),

		Functions: make(map[string]int64),
		fnIndex:   0,

		Symbols:     make(map[string]int64),
		symbolIndex: 0,

		floats:     make(map[int64]float64),
		floatIndex: 0,
	}
}

func (env *Env) Push(e int64) {
	if len(env.Stack) == cap(env.Stack) {
		bigStack := make([]int64, len(env.Stack), cap(env.Stack)*2)
		copy(bigStack, env.Stack)
		env.Stack = bigStack
	}

	env.Stack = append(env.Stack, e)
}

func (env *Env) Pop() int64 {
	e := env.Stack[len(env.Stack)-1]
	env.Stack = env.Stack[:len(env.Stack)-1]
	return e
}

func (env *Env) PushFloat(f float64) {
	env.floats[env.floatIndex] = f
	env.Push(env.floatIndex)
	env.floatIndex++
}

func (env *Env) GetFloat(index int64) float64 {
	return env.floats[index]
}
