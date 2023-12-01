package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jimschubert/answer/input"
)

func main() {
	m := input.New()
	m.Prompt = "Please enter your name:"
	m.Placeholder = "(first name only)"
	m.Validate = func(v string) error {
		if v == "" {
			return nil
		}
		if len(v) >= 2 && !unicode.IsUpper(rune(v[0])) {
			return errors.New("Name must be uppercase")
		}
		return nil
	}
	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Hi, %s!\n", m.Value())
}
