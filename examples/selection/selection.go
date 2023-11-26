package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jimschubert/answer/colors"
	"github.com/jimschubert/answer/selection"
)

func main() {
	m := selection.New()
	m.Prompt = "Please select your three favorite letters:"
	m.MaxSelections = 3
	m.ChooserIndicator = '✎'
	m.Styles.Text = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: colors.TextLight, Dark: colors.TextDark})
	choices := make([]string, 0)
	for i := 'A'; i <= 'Z'; i++ {
		choices = append(choices, string(i))
	}
	m.Choices = choices
	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "You selected: %v\n", m.SelectedValues())
}
