package main

type Variable struct {
	Value int64
	Type  YapType
}

type Function struct {
	Argc int64
	Body Stack
}

type Env struct {
	functions map[int64]Function
	fnIndex   int64

	symbols       map[string]int64
	variables     map[int64]Variable
	variableIndex int64

	floats     map[int64]float64
	floatIndex int64
}

func NewEnv() Env {
	return Env{
		functions: make(map[int64]Function),
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

func (env *Env) CreateFn(argc int64, body Stack) int64 {
	flipStack(&body)
	env.functions[env.fnIndex] = Function{
		argc,
		body,
	}
	env.fnIndex++
	return env.fnIndex - 1
}

func (env *Env) GetFn(id int64) (int64, Stack) {
	if f, ok := env.functions[id]; ok {
		return f.Argc, f.Body
	}

	panic("Could not find function")
}
