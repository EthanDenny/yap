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

		for len(env.Stack) > 1 {
			eval(env)
		}

		if len(env.Stack) > 1 {
			fmt.Println(env.Stack[0])
		}
	}
}
