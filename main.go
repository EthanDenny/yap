package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	env := NewEnv()
	symbols := NewSymbolTable(nil)

	if len(os.Args) >= 2 {
		fileName := os.Args[1]
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		var stack Stack
		var tokenList TokenList

		fileScanner := bufio.NewScanner(file)
		for fileScanner.Scan() {
			line := fileScanner.Text()
			line = strings.TrimSpace(line)
			tokenList.Append(scan(line))
		}

		parse(env, symbols, &stack, tokenList)

		for len(stack) > 0 {
			eval(env, symbols, &stack)
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

			parse(env, symbols, &stack, tokenList)
			eval(env, symbols, &stack)
		}
	}
}
