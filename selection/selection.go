package selection

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jimschubert/answer/colors"
)

var (
	_ tea.Model = (*Model)(nil)
)

// Styles holds relevant styles used for rendering
type Styles struct {
	PromptPrefix      lipgloss.Style
	Prompt            lipgloss.Style
	Text              lipgloss.Style
	SelectedIndicator lipgloss.Style
	ChooserIndicator  lipgloss.Style
}

// Model represents the bubble tea model for the selection
type Model struct {
	PromptPrefix      string
	Prompt            string
	SelectedIndicator rune
	ChooserIndicator  rune
	Styles            Styles
	Choices           []string
	KeyMap            KeyMap
	MaxSelections     int
	HideHelp          bool
	PerPage           int
	cursor            int
	paginator         paginator.Model
	help              help.Model
	initialized       bool
	selected          map[int]struct{}
}

type KeyMap struct {
	SelectionUp   key.Binding
	SelectionDown key.Binding
	PageNext      key.Binding
	PagePrev      key.Binding
	Quit          key.Binding
	Select        key.Binding
	Help          key.Binding
	Enter         key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SelectionUp, k.SelectionDown, k.Select},
		{k.PagePrev, k.PageNext},
		{k.Help, k.Quit},
	}
}

var DefaultKeyMap = KeyMap{
	SelectionUp: key.NewBinding(
		key.WithKeys("k", tea.KeyUp.String()),
		key.WithHelp("↑/k", "up"),
	),
	SelectionDown: key.NewBinding(
		key.WithKeys("j", tea.KeyDown.String()),
		key.WithHelp("↓/j", "down"),
	),
	PagePrev: key.NewBinding(
		key.WithKeys("h", tea.KeyLeft.String()),
		key.WithHelp("←/h", "prev"),
	),
	PageNext: key.NewBinding(
		key.WithKeys("l", tea.KeyRight.String()),
		key.WithHelp("→/l", "next"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", tea.KeyEsc.String(), tea.KeyCtrlC.String()),
		key.WithHelp("q", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys(tea.KeySpace.String()),
		key.WithHelp("space", "select"),
	),
	Enter: key.NewBinding(key.WithKeys(tea.KeyEnter.String())),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

// New creates a new model with default settings.
func New() Model {
	return Model{
		PromptPrefix:      "? ",
		KeyMap:            DefaultKeyMap,
		SelectedIndicator: 'x',
		ChooserIndicator:  '➤',
		Styles: Styles{
			PromptPrefix:      lipgloss.NewStyle().Foreground(lipgloss.Color(colors.PromptPrefix)),
			SelectedIndicator: lipgloss.NewStyle().Foreground(lipgloss.Color(colors.PromptPrefix)),
			ChooserIndicator:  lipgloss.NewStyle().Foreground(lipgloss.Color(colors.PromptPrefix)),
		},
		help:     help.New(),
		selected: make(map[int]struct{}),
	}
}

func (m *Model) setup() {
	if m.Prompt == "" {
		m.Prompt = "Please select:"
	}

	paginate := paginator.New()
	paginate.Type = paginator.Dots
	if m.PerPage < 1 {
		paginate.PerPage = 10
	} else {
		paginate.PerPage = m.PerPage
	}
	paginate.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	paginate.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	paginate.KeyMap.NextPage = m.KeyMap.PageNext
	paginate.KeyMap.PrevPage = m.KeyMap.PagePrev
	paginate.SetTotalPages(len(m.Choices))

	m.paginator = paginate
	m.initialized = true
}

func (m *Model) Init() tea.Cmd {
	m.setup()
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.paginator, cmd = m.paginator.Update(msg)
	start, _ := m.paginator.GetSliceBounds(len(m.Choices))

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate its view as needed.
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Quit, m.KeyMap.Enter):
			m.HideHelp = true
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.SelectionUp):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.KeyMap.SelectionDown):
			if m.cursor < len(m.Choices)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.KeyMap.PageNext, m.KeyMap.PagePrev):
			m.cursor = 0
		case key.Matches(msg, m.KeyMap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.KeyMap.Select):
			idx := start + m.cursor
			if _, ok := m.selected[idx]; ok {
				delete(m.selected, idx)
			} else {
				m.selected[idx] = struct{}{}
			}
			if m.MaxSelections == 1 {
				return m, tea.Quit
			}
		}
	}

	return m, cmd
}

func (m *Model) SelectedIndexes() []int {
	indexes := make([]int, 0)
	for idx := range m.selected {
		indexes = append(indexes, idx)
	}
	sort.Ints(indexes)
	return indexes
}

func (m *Model) SelectedValues() []string {
	values := make([]string, 0)
	for i, choice := range m.Choices {
		if _, ok := m.selected[i]; ok {
			values = append(values, choice)
		}
	}
	return values
}

func (m *Model) View() string {

	styleText := m.Styles.Text.Inline(true).Render

	var b strings.Builder
	b.WriteString(m.Styles.PromptPrefix.Inline(true).Render(m.PromptPrefix))
	if !strings.HasSuffix(m.PromptPrefix, " ") {
		b.WriteString(" ")
	}
	b.WriteString(m.Styles.Prompt.Render(m.Prompt))
	b.WriteString("\n\n")

	start, end := m.paginator.GetSliceBounds(len(m.Choices))
	for i, item := range m.Choices[start:end] {
		idx := start + i
		cursor := " "
		if m.cursor == i {
			cursor = m.Styles.ChooserIndicator.Inline(true).Render(string(m.ChooserIndicator))
		}

		b.WriteString(cursor)
		b.WriteString(" ")
		if m.MaxSelections != 1 {
			b.WriteString(styleText("["))
			if _, ok := m.selected[idx]; ok {
				b.WriteString(styleText(m.Styles.SelectedIndicator.Inline(true).Render(string(m.SelectedIndicator))))
			} else {
				b.WriteString(" ")
			}
			b.WriteString(styleText("] "))
		}
		b.WriteString(styleText(item))
		b.WriteString("\n")
	}
	b.WriteString("  " + m.paginator.View())
	if !m.HideHelp {
		helpView := m.help.View(m.KeyMap)
		b.WriteString("\n\n")
		b.WriteString(helpView)
	}
	b.WriteString("\n")
	return b.String()
}
