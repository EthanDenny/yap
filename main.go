package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ethandenny/yap/tokens"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	if len(os.Args) >= 2 {
		fileName := os.Args[1]
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		env := NewEnv()
		var stack Stack
		var tokenList tokens.TokenList

		fileScanner := bufio.NewScanner(file)
		for fileScanner.Scan() {
			line := fileScanner.Text()
			line = strings.TrimSpace(line)
			tokenList.Append(scan(line))
		}

		parse(&env, &stack, tokenList)

		for len(stack) > 0 {
			eval(&env, &stack, nil)
		}
	} else {
		env := NewEnv()
		for {
			fmt.Print("> ")
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			fmt.Println()

			line = strings.TrimSpace(line)
			tokenList := scan(line)

			var stack Stack

			parse(&env, &stack, tokenList)
			eval(&env, &stack, nil)
		}
	}
}
