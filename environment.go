package main

type SymbolTable struct {
	parent *SymbolTable
	table  map[string]int64
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		parent: parent,
		table:  make(map[string]int64),
	}
}

func (s *SymbolTable) Get(name string) int64 {
	if index, ok := s.table[name]; ok {
		return index
	} else if s.parent != nil {
		return s.parent.Get(name)
	} else {
		return -1
	}
}

type Variable struct {
	Value int64
	Type  YapType
}

type Function struct {
	ArgNames []string
	Body     TokenList
}

type Env struct {
	symbols map[string]int64

	functions map[int64]Function
	fnIndex   int64

	variables     map[int64]Variable
	variableIndex int64

	floats     map[int64]float64
	floatIndex int64

	strings     map[int64]string
	stringIndex int64
}

func NewEnv() *Env {
	return &Env{
		functions: make(map[int64]Function),
		fnIndex:   0,

		variables:     make(map[int64]Variable),
		variableIndex: 0,

		floats:     make(map[int64]float64),
		floatIndex: 0,

		strings:     make(map[int64]string),
		stringIndex: 0,
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

func (env *Env) InsertString(s string) int64 {
	env.strings[env.stringIndex] = s
	index := env.stringIndex
	env.stringIndex++
	return index
}

func (env *Env) GetString(index int64) string {
	if str, ok := env.strings[index]; ok {
		return str
	}

	panic("Could not find string")
}

func (env *Env) SetVariable(symbols *SymbolTable, name string, value int64, Type YapType) {
	if id, ok := symbols.table[name]; ok {
		env.variables[id] = Variable{
			value,
			Type,
		}
	} else {
		symbols.table[name] = env.variableIndex
		env.variables[env.variableIndex] = Variable{
			value,
			Type,
		}
		env.variableIndex++
	}
}

func (env *Env) GetVariable(id int64) (int64, YapType) {
	if v, ok := env.variables[id]; ok {
		return v.Value, v.Type
	}

	panic("Could not find variable")
}

func (env *Env) CreateFn(argNames []string) int64 {
	env.functions[env.fnIndex] = Function{
		argNames,
		TokenList{},
	}
	env.fnIndex++
	return env.fnIndex - 1
}

func (env *Env) SetFnBody(id int64, body TokenList) {
	fn := env.functions[id]
	fn.Body = body
	env.functions[id] = fn
}

func (env *Env) GetFn(id int64) ([]string, TokenList) {
	if f, ok := env.functions[id]; ok {
		return f.ArgNames, f.Body
	}

	panic("Could not find function")
}
