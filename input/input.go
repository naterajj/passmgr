package input

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func Read(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(errors.New("An error occurred reading input"))
	}
	return strings.Trim(input, "\n")
}

func ReadSecret(prompt string) (string, error) {
	fmt.Print(prompt + ": ")
	secret1, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(errors.New("Unable to read secret from terminal."))
	}
	fmt.Print("\n")
	fmt.Print(prompt + " (again): ")
	secret2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(errors.New("Unable to read secret from terminal."))
	}
	fmt.Print("\n")
	if string(secret1) != string(secret2) {
		return "", errors.New("Secrets mismatch, try again")
	}
	return string(secret1), nil
}
