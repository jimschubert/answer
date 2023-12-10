package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jimschubert/answer/input"
	"github.com/jimschubert/answer/suggest"
)

func main() {
	m := input.New()
	m.Prompt = "Please enter your name:"
	m.Placeholder = "(first name only)"
	m.Suggest = suggest.LevenshteinDistance([]string{"Jim", "James", "Jameson"},
		suggest.LevenshteinDistanceMin(0),
		suggest.LevenshteinDistanceMax(4))
	m.Validate = func(v string) error {
		if v == "" {
			return nil
		}
		if !unicode.IsUpper(rune(v[0])) {
			return errors.New("name must be uppercase")
		}
		return nil
	}
	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Hi, %s!\n", m.Value())
}
