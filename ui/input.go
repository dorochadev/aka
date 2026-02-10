package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Prompt asks a question and returns the user's input
func Prompt(question string) string {
	CurrentTheme.Primary.Printf("%s ", question)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(answer)
}

// PromptDefault asks a question with a default value
func PromptDefault(question, defaultValue string) string {
	CurrentTheme.Primary.Printf("%s ", question)
	CurrentTheme.Placeholder.Printf("[%s] ", defaultValue)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return defaultValue
	}
	return answer
}

// Confirm asks a yes/no question and returns a boolean
func Confirm(question string) bool {
	CurrentTheme.Primary.Printf("%s ", question)
	CurrentTheme.Muted.Print("(y/n) ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}

// PromptPassword securely prompts for a password without echoing
func PromptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(password), nil
}
