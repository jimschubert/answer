package input

import (
	"bytes"
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
	buf := bytes.Buffer{}
	if m.PromptPrefix != "" {
		buf.WriteString(m.Styles.PromptPrefix.Render(m.PromptPrefix))
		if m.Prompt != "" && !strings.HasSuffix(m.PromptPrefix, " ") {
			buf.WriteRune(' ')
		}
	}

	if m.done {
		if m.Prompt != "" {
			buf.WriteString(m.Styles.Prompt.Render(m.Prompt))
			buf.WriteRune(' ')
		}
		buf.WriteString(m.input.Value())
		buf.WriteRune('\n')
		return buf.String()
	}

	buf.WriteString(m.input.View())
	if m.err != nil {
		buf.WriteRune('\n')
		buf.WriteString(m.Styles.ErrorPrefix.Render("✘"))
		buf.WriteString(m.Styles.Placeholder.Render(": " + m.err.Error() + "\n"))
	}
	return buf.String()
}
