package main

import (
	"strconv"

	"github.com/ethandenny/yap/tokens"
)

func parse(list tokens.TokenList) *Env {
	env := NewEnv()
	var instructions []int64

	for list.HasToken() {
		for _, instr := range parseCall(&env, &list) {
			instructions = append(instructions, instr)
		}
	}

	for i := len(instructions) - 1; i >= 0; i-- {
		env.Push(instructions[i])
	}

	return &env
}

func parseCall(env *Env, list *tokens.TokenList) []int64 {
	list.Expect(tokens.LeftParen)

	callName := list.Expect(tokens.Symbol).Content

	var args [][]int64

	for nextToken := list.Peek(); nextToken.Type != tokens.RightParen; nextToken = list.Peek() {
		switch nextToken.Type {
		case tokens.Integer:
			t := list.Consume()
			i, _ := strconv.ParseInt(t.Content, 10, 64)
			args = append(args, []int64{Integer, i})
		case tokens.Float:
			t := list.Consume()
			f, _ := strconv.ParseFloat(t.Content, 64)
			index := env.InsertFloat(f)
			args = append(args, []int64{Float, index})
		case tokens.Symbol:
			symbolIndex := env.GetSymbol(nextToken.Content)
			args = append(args, []int64{Var, symbolIndex})
		case tokens.LeftParen:
			args = append(args, parseCall(env, list))
		default:
		}
	}

	var functions = map[string]int64{
		"+":   Add,
		"yap": Print,
	}

	var instructions []int64

	if f, containsKey := functions[callName]; containsKey {
		instructions = append(instructions, f)
	} else {
		fnIndex := env.GetFn(callName)
		instructions = append(instructions, fnIndex)
	}

	instructions = append(instructions, int64(len(args)))

	for _, arg := range args {
		for _, instr := range arg {
			instructions = append(instructions, instr)
		}
	}

	list.Expect(tokens.RightParen)

	return instructions
}
