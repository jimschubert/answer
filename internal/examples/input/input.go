package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jimschubert/answer/input"
	"github.com/jimschubert/answer/suggest"
	"github.com/jimschubert/answer/validate"
)

var complexValidation = flag.Bool("complex", false, "use complex validations")

func main() {
	flag.Parse()

	m := input.New()
	m.Prompt = "Please enter your name:"
	m.Placeholder = "(first name only)"
	m.Suggest = suggest.LevenshteinDistance([]string{"Jim", "James", "Jameson"},
		suggest.LevenshteinDistanceMin(0),
		suggest.LevenshteinDistanceMax(4))

	requireUppercase := func(v string) error {
		if v != "" && !unicode.IsUpper(rune(v[0])) {
			return errors.New("name must be uppercase")
		}
		return nil
	}

	if *complexValidation {
		m.Validate = validate.NewValidation().
			MinLength(2, "min: 2 characters").
			MaxLength(5, "max: 8 characters").
			And(func(input string) error {
				for _, r := range input {
					if !unicode.IsLetter(r) {
						return errors.New("letters only")
					}
				}
				return nil
			}).
			And(requireUppercase).
			Build()
	} else {
		m.Validate = requireUppercase
	}

	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	result := m.Value()
	if result != "" {
		_, _ = fmt.Fprintf(os.Stdout, "Hi, %s!\n", m.Value())
	}
}
