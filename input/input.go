package input

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jimschubert/answer/colors"
)

var (
	_ tea.Model = (*Model)(nil)
)

// ValidateFunc determines if the input string is valid, returning nil if valid or an error if invalid
type ValidateFunc func(input string) error

// Styles holds relevant styles used for rendering
type Styles struct {
	PromptPrefix lipgloss.Style
	Prompt       lipgloss.Style
	ErrorPrefix  lipgloss.Style
	Text         lipgloss.Style
	Placeholder  lipgloss.Style
}

// Model represents the bubble tea model for the input
type Model struct {
	PromptPrefix string
	Prompt       string
	Placeholder  string
	CharLimit    int
	MaxWidth     int
	EchoMode     textinput.EchoMode
	Validate     ValidateFunc
	Styles       Styles
	err          error
	done         bool
	input        textinput.Model
	initialized  bool
}

// New creates a new model with default settings.
func New() Model {
	return Model{
		PromptPrefix: "? ",
		CharLimit:    0,
		Styles: Styles{
			PromptPrefix: lipgloss.NewStyle().Foreground(lipgloss.Color(colors.PromptPrefix)),
			ErrorPrefix:  lipgloss.NewStyle().Foreground(lipgloss.Color(colors.ErrorPrefix)),
			Placeholder:  lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Placeholder)),
		},
	}
}

func (m *Model) setup() {
	if m.Validate == nil {
		m.Validate = func(input string) error {
			return nil
		}
	}
	if m.Prompt == "" {
		m.Prompt = "Please enter:"
	}
	input := textinput.New()
	input.CharLimit = m.CharLimit
	input.Width = m.MaxWidth
	if !strings.HasSuffix(m.Prompt, " ") {
		input.Prompt = m.Prompt + " "
	} else {
		input.Prompt = m.Prompt
	}
	input.Placeholder = m.Placeholder
	input.PromptStyle = m.Styles.Prompt
	input.PlaceholderStyle = m.Styles.Placeholder
	input.TextStyle = m.Styles.Text
	input.EchoMode = m.EchoMode
	input.Focus()
	m.input = input
	m.initialized = true
}

func (m *Model) Init() tea.Cmd {
	m.setup()
	return nil
}

func (m *Model) SetValue(value string) {
	m.input.SetValue(value)
}

func (m *Model) Value() string {
	return m.input.Value()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.initialized {
		m.setup()
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyEnter:
			if m.err == nil {
				m.done = true
				return m, tea.Quit
			}
		}
	case error:
		m.err = msg
	}

	var cmds []tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.err = m.Validate(m.input.Value())

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var b strings.Builder
	if m.PromptPrefix != "" {
		b.WriteString(m.Styles.PromptPrefix.Inline(true).Render(m.PromptPrefix))
		if m.Prompt != "" && !strings.HasSuffix(m.PromptPrefix, " ") {
			b.WriteRune(' ')
		}
	}

	if m.done {
		// rather than clearing the program output, we want to show the question + answer just as AlecAivazis/survey did
		if m.Prompt != "" {
			b.WriteString(m.Styles.Prompt.Inline(true).Render(m.Prompt))
			b.WriteRune(' ')
		}
		b.WriteString(m.input.Value())
		b.WriteRune('\n')
		return b.String()
	}

	b.WriteString(m.input.View())
	if m.err != nil {
		b.WriteRune('\n')
		b.WriteString(m.Styles.ErrorPrefix.Inline(true).Render("✘"))
		b.WriteString(m.Styles.Placeholder.Inline(true).Render(": " + m.err.Error() + "\n"))
	}
	return b.String()
}
