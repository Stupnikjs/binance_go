package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt() string {

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input:", err)
		return ""
	}

	// The ReadString function includes the newline character, so we need to trim it.
	command := strings.TrimSpace(input)

	return command
}
