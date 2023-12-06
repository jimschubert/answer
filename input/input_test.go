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
	"github.com/jimschubert/stripansi"
	"github.com/stretchr/testify/assert"
)

type state struct {
	Name        string
	BeforeType  string
	Inputs      []tea.KeyMsg
	AfterType   string
	ExpectView  string
	ExpectValue string
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
					ExpectView:  "? Please enter your name: jim \r\n✘ Name must be uppercase",
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
					tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
				}

				time.Sleep(100 * time.Millisecond)
				if err := tm.Quit(); err != nil {
					t.Fatal(err)
				}

				out := stripansi.Bytes(readBts(t, tm.FinalOutput(t, teatest.WithFinalTimeout(2*time.Second))))

				if s.ExpectView != "" {
					actual := stripansi.Bytes(out)
					assert.Contains(t, string(actual), s.ExpectView)
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
