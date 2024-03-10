package main

type Env struct {
	Stack []int64

	functions map[string]int64
	fnIndex   int64

	symbols     map[string]int64
	symbolIndex int64

	floats     map[int64]float64
	floatIndex int64
}

func NewEnv() Env {
	return Env{
		Stack: make([]int64, 0, 1),

		functions: make(map[string]int64),
		fnIndex:   0,

		symbols:     make(map[string]int64),
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

func (env *Env) InsertFloat(f float64) int64 {
	env.floats[env.floatIndex] = f
	index := env.floatIndex
	env.floatIndex++
	return index
}

func (env *Env) GetFloat(index int64) float64 {
	return env.floats[index]
}

func (env *Env) GetSymbol(name string) int64 {
	return env.symbols[name]
}

func (env *Env) GetFn(name string) int64 {
	return env.functions[name]
}
