package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)

		tokenList := scan(line)

		env := parse(tokenList)

		for stackPos := len(env.Stack) - 1; stackPos >= 0; stackPos-- {
			dis(env, &stackPos)
		}

		fmt.Println()

		for len(env.Stack) > 1 {
			eval(env)
		}
	}
}

func dis(env *Env, stackPos *int) {
	switch env.Stack[*stackPos] {
	case Add:
		fmt.Println("ADD")
		argc := getArgc(env, stackPos)
		fmt.Println("ARGC:", argc)
		disArg(env, stackPos)
		disArg(env, stackPos)
	case Print:
		fmt.Println("PRINT")
		argc := getArgc(env, stackPos)
		fmt.Println("ARGC:", argc)
		var i int64 = 0
		for ; i < argc; i++ {
			disArg(env, stackPos)
		}
	}

}

func getArgc(env *Env, stackPos *int) int64 {
	*stackPos--
	argc := env.Stack[*stackPos]
	*stackPos--
	return argc
}

func disArg(env *Env, stackPos *int) {
	switch env.Stack[*stackPos] {
	case Integer:
		fmt.Println("INTEGER")
		*stackPos--
		fmt.Println(env.Stack[*stackPos])
		*stackPos--
	case Float:
		fmt.Println("FLOAT")
		*stackPos--
		f := env.GetFloat(env.Stack[*stackPos])
		fmt.Println(f)
		*stackPos--
	case Add:
		dis(env, stackPos)
	case Print:
		dis(env, stackPos)
	}
}
