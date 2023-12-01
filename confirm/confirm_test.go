package confirm

import (
	"fmt"
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/jimschubert/stripansi"
	"github.com/stretchr/testify/assert"
)

type state struct {
	Name        string
	Inputs      []tea.KeyMsg
	ExpectView  string
	ExpectValue Decision
}

func TestModel_View(t *testing.T) {
	tests := []struct {
		name   string
		model  Model
		states []state
	}{
		{
			name: "input style",
			model: func() Model {
				m := New()
				m.Prompt = "Do you like testing?"
				return m
			}(),
			states: []state{
				{
					Name:        "selects the default",
					Inputs:      []tea.KeyMsg{},
					ExpectView:  "? Do you like testing? y\r\n",
					ExpectValue: Accepted,
				},
				{
					Name:        "selects via lowercase y",
					Inputs:      []tea.KeyMsg{{Runes: []rune{'y'}}},
					ExpectView:  "? Do you like testing? y\r\n",
					ExpectValue: Accepted,
				},
				{
					Name:        "selects via uppercase Y",
					Inputs:      []tea.KeyMsg{{Runes: []rune{'Y'}}},
					ExpectView:  "? Do you like testing? y\r\n",
					ExpectValue: Accepted,
				},
				{
					Name:        "selects via lowercase n",
					Inputs:      []tea.KeyMsg{{Runes: []rune{'n'}}},
					ExpectView:  "? Do you like testing? n\r\n",
					ExpectValue: Denied,
				},
				{
					Name:        "selects via uppercase N",
					Inputs:      []tea.KeyMsg{{Runes: []rune{'N'}}},
					ExpectView:  "? Do you like testing? n\r\n",
					ExpectValue: Denied,
				},
				{
					Name:        "disallows invalid input",
					Inputs:      []tea.KeyMsg{{Runes: []rune{'R'}}},
					ExpectView:  "? Do you like testing? y\r\n", // accepts the default
					ExpectValue: Accepted,                       // the default
				},
			},
		},
		{
			name: "input style with undecided default",
			model: func() Model {
				m := New()
				m.DefaultValue = Undecided
				m.Prompt = "Do you like testing?"
				return m
			}(),
			states: []state{
				{
					Name:        "selects the default",
					Inputs:      []tea.KeyMsg{},
					ExpectView:  "? Do you like testing? \r\n",
					ExpectValue: Undecided,
				},
				{
					Name:        "selects via lowercase y",
					Inputs:      []tea.KeyMsg{{Runes: []rune{'y'}}},
					ExpectView:  "? Do you like testing? y\r\n",
					ExpectValue: Accepted,
				},
				{
					Name:        "selects via uppercase Y",
					Inputs:      []tea.KeyMsg{{Runes: []rune{'Y'}}},
					ExpectView:  "? Do you like testing? y\r\n",
					ExpectValue: Accepted,
				},
				{
					Name:        "selects via lowercase n",
					Inputs:      []tea.KeyMsg{{Runes: []rune{'n'}}},
					ExpectView:  "? Do you like testing? n\r\n",
					ExpectValue: Denied,
				},
				{
					Name:        "selects via uppercase N",
					Inputs:      []tea.KeyMsg{{Runes: []rune{'N'}}},
					ExpectView:  "? Do you like testing? n\r\n",
					ExpectValue: Denied,
				},
			},
		},

		{
			name: "horizontal style",
			model: func() Model {
				m := New()
				m.Prompt = "Do you like testing?"
				m.Rendering = HorizontalSelection
				return m
			}(),
			states: []state{
				{
					Name:        "selects the default",
					Inputs:      []tea.KeyMsg{},
					ExpectView:  "➤y  n",
					ExpectValue: Accepted,
				},
				{
					Name:        "selects via right toggle",
					Inputs:      []tea.KeyMsg{{Type: tea.KeyRight}},
					ExpectView:  " y ➤n",
					ExpectValue: Denied,
				},
				{
					Name:        "selects via right toggle cycling",
					Inputs:      []tea.KeyMsg{{Type: tea.KeyRight}, {Type: tea.KeyRight}},
					ExpectView:  "➤y  n",
					ExpectValue: Accepted,
				},
				{
					Name:        "selects via left toggle",
					Inputs:      []tea.KeyMsg{{Type: tea.KeyLeft}},
					ExpectView:  " y ➤n",
					ExpectValue: Denied,
				},
				{
					Name:        "selects via left toggle cycling",
					Inputs:      []tea.KeyMsg{{Type: tea.KeyLeft}, {Type: tea.KeyLeft}},
					ExpectView:  "➤y  n",
					ExpectValue: Accepted,
				},
			},
		},
		{
			name: "vertical style",
			model: func() Model {
				m := New()
				m.Prompt = "Do you like testing?"
				m.Rendering = VerticalSelection
				return m
			}(),
			states: []state{
				{
					Name:        "selects the default",
					Inputs:      []tea.KeyMsg{},
					ExpectView:  "➤ y",
					ExpectValue: Accepted,
				},
				{
					Name:        "selects via up toggle",
					Inputs:      []tea.KeyMsg{{Type: tea.KeyUp}},
					ExpectView:  "➤ n",
					ExpectValue: Denied,
				},
				{
					Name:        "selects via up toggle cycling",
					Inputs:      []tea.KeyMsg{{Type: tea.KeyUp}, {Type: tea.KeyUp}},
					ExpectView:  "➤ y",
					ExpectValue: Accepted,
				},
				{
					Name:        "selects via down toggle",
					Inputs:      []tea.KeyMsg{{Type: tea.KeyDown}},
					ExpectView:  "➤ n",
					ExpectValue: Denied,
				},
				{
					Name:        "selects via down toggle cycling",
					Inputs:      []tea.KeyMsg{{Type: tea.KeyDown}, {Type: tea.KeyDown}},
					ExpectView:  "➤ y",
					ExpectValue: Accepted,
				},
			},
		},
	}
	for _, tt := range tests {
		for _, s := range tt.states {
			t.Run(fmt.Sprintf("%s_%s", tt.name, s.Name), func(t *testing.T) {
				tm := teatest.NewTestModel(
					t, &tt.model,
					teatest.WithInitialTermSize(120, 40),
				)

				t.Cleanup(func() {
					if err := tm.Quit(); err != nil {
						t.Fatal(err)
					}
				})

				for _, input := range s.Inputs {
					tm.Send(input)
					time.Sleep(200 * time.Millisecond)
				}

				tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
				time.Sleep(100 * time.Millisecond)
				if err := tm.Quit(); err != nil {
					t.Fatal(err)
				}

				out := stripansi.Bytes(readBts(t, tm.FinalOutput(t, teatest.WithFinalTimeout(2*time.Second))))
				if s.ExpectView != "" {
					actual := stripansi.Bytes(out)
					assert.Contains(t, string(actual), s.ExpectView)
				}

				assert.Equal(t, s.ExpectValue, tt.model.Selected())
				teatest.RequireEqualOutput(t, out)
			})
		}
	}
}

func readBts(tb testing.TB, r io.Reader) []byte {
	tb.Helper()
	bts, err := io.ReadAll(r)
	if err != nil {
		tb.Fatal(err)
	}
	return bts
}
