package main

type Variable struct {
	Value int64
	Type  YapType
}

type Env struct {
	functions map[string]int64
	fnIndex   int64

	symbols       map[string]int64
	variables     map[int64]Variable
	variableIndex int64

	floats     map[int64]float64
	floatIndex int64
}

func NewEnv() Env {
	return Env{
		functions: make(map[string]int64),
		fnIndex:   0,

		symbols:       make(map[string]int64),
		variables:     make(map[int64]Variable),
		variableIndex: 0,

		floats:     make(map[int64]float64),
		floatIndex: 0,
	}
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

func (env *Env) SetVariable(name string, value int64, Type YapType) {
	var index int64

	if id, ok := env.symbols[name]; ok {
		index = id
	} else {
		index = env.variableIndex
		env.symbols[name] = index
		env.variableIndex++
	}

	env.variables[index] = Variable{
		value,
		Type,
	}
}

func (env *Env) GetSymbol(name string) int64 {
	if id, ok := env.symbols[name]; ok {
		return id
	}

	panic("Could not find symbol")
}

func (env *Env) GetVariable(id int64) (int64, YapType) {
	if v, ok := env.variables[id]; ok {
		return v.Value, v.Type
	}

	panic("Could not find variable")
}

func (env *Env) GetFn(name string) int64 {
	if f, ok := env.functions[name]; ok {
		return f
	}

	panic("Could not find function")
}
