package input

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/jimschubert/answer/suggest"
	"github.com/jimschubert/stripansi"
	"github.com/stretchr/testify/assert"
)

type state struct {
	Name                     string
	BeforeType               string
	Inputs                   []tea.KeyMsg
	AfterType                string
	ExpectView               string
	ExpectViewDoesNotContain string
	ExpectValue              string
}

func TestModel_View(t *testing.T) {
	tests := []struct {
		name       string
		inputModel Model
		states     []state
	}{
		{
			name: "validatable input",
			inputModel: func() Model {
				m := New()
				m.Prompt = "Please enter your name:"
				m.Placeholder = "(first name only)"
				m.Validate = func(v string) error {
					if v == "" {
						return nil
					}
					if len(v) >= 2 && !unicode.IsUpper(rune(v[0])) {
						return errors.New("name must be uppercase")
					}
					return nil
				}
				return m
			}(),
			states: []state{
				{
					Name:        "displays validation message",
					BeforeType:  "jim",
					Inputs:      []tea.KeyMsg{},
					ExpectView:  "? Please enter your name: jim \r\n✘ name must be uppercase",
					ExpectValue: "jim",
				},
				{
					Name:        "removes validation message after fix",
					BeforeType:  "jim",
					Inputs:      []tea.KeyMsg{{Type: tea.KeyBackspace}, {Type: tea.KeyBackspace}, {Type: tea.KeyBackspace}},
					AfterType:   "Jim",
					ExpectView:  "? Please enter your name: Jim",
					ExpectValue: "Jim",
				},
			},
		},

		{
			name: "input with suggestions does not persist suggestions on enter",
			inputModel: func() Model {
				m := New()
				m.Prompt = "Please enter your name:"
				m.Placeholder = "(first name only)"
				m.Validate = func(v string) error {
					if v == "" {
						return nil
					}
					if len(v) >= 2 && !unicode.IsUpper(rune(v[0])) {
						return errors.New("name must be uppercase")
					}
					return nil
				}

				m.Suggest = suggest.LevenshteinDistance([]string{"Jim", "James", "Jameson"},
					suggest.LevenshteinDistanceMin(0),
					suggest.LevenshteinDistanceMax(4))

				return m
			}(),
			states: []state{
				{
					Name:                     "when suggestions are not displayed",
					BeforeType:               "Ji",
					Inputs:                   []tea.KeyMsg{},
					ExpectValue:              "Ji",
					ExpectViewDoesNotContain: "Suggestions:",
				},
				{
					Name:                     "when suggestions are displayed",
					BeforeType:               "Jim",
					Inputs:                   []tea.KeyMsg{{Type: tea.KeyBackspace}, {Type: tea.KeyBackspace}},
					AfterType:                "ames",
					ExpectValue:              "James",
					ExpectViewDoesNotContain: "Suggestions:",
				},
			},
		},
	}

	for _, tt := range tests {
		for _, s := range tt.states {
			t.Run(fmt.Sprintf("%s_%s", tt.name, s.Name), func(t *testing.T) {

				tm := teatest.NewTestModel(
					t, &tt.inputModel,
					teatest.WithInitialTermSize(120, 40),
				)

				t.Cleanup(func() {
					if err := tm.Quit(); err != nil {
						t.Fatal(err)
					}
				})

				if s.BeforeType != "" {
					tm.Type(s.BeforeType)
				}
				for _, input := range s.Inputs {
					tm.Send(input)
					time.Sleep(100 * time.Millisecond)
				}

				if s.AfterType == "" {
					tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
				} else {
					tm.Type(s.AfterType)
				}

				time.Sleep(100 * time.Millisecond)
				if err := tm.Quit(); err != nil {
					t.Fatal(err)
				}

				out := stripansi.Bytes(readBts(t, tm.FinalOutput(t, teatest.WithFinalTimeout(2*time.Second))))
				actual := stripansi.Bytes(out)

				if s.ExpectView != "" {
					assert.Contains(t, string(actual), s.ExpectView)
				}

				if s.ExpectViewDoesNotContain != "" {
					assert.NotContainsf(t, string(actual), s.ExpectViewDoesNotContain, "View should not contain %s", s.ExpectViewDoesNotContain)
				}

				model := tm.FinalModel(t).(*Model)
				assert.Equal(t, s.ExpectValue, model.Value())

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
