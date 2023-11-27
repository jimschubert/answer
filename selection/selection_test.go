package selection

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
	Name                  string
	Inputs                []tea.KeyMsg
	ExpectView            string
	ExpectSelectedIndexes []int
	ExpectSelectedValues  []string
}

func TestModel_View(t *testing.T) {
	tests := []struct {
		name           string
		selectionModel Model
		states         []state
	}{
		{
			name: "single page display",
			selectionModel: func() Model {
				m := New()
				m.Prompt = "Choose a color:"
				m.Choices = []string{"Red", "Green", "Blue"}
				return m
			}(),
			states: []state{
				{
					Name:                  "selects item",
					Inputs:                []tea.KeyMsg{{Type: tea.KeyDown}, {Type: tea.KeySpace}, {Type: tea.KeyEnter}},
					ExpectView:            "? Choose a color:",
					ExpectSelectedIndexes: []int{1},
					ExpectSelectedValues:  []string{"Green"},
				},
			},
		},
		{
			name: "multi page display",
			selectionModel: func() Model {
				m := New()
				m.Prompt = "Please select your favorite letters:"
				m.PerPage = 6
				choices := make([]string, 0)
				for i := 'A'; i <= 'Z'; i++ {
					choices = append(choices, string(i))
				}
				m.Choices = choices
				return m
			}(),
			states: []state{
				{
					Name: "selects multiple",
					Inputs: []tea.KeyMsg{
						{Type: tea.KeyDown},
						{Type: tea.KeyDown},
						{Type: tea.KeySpace},
						{Type: tea.KeyDown},
						{Type: tea.KeySpace},
						{Type: tea.KeyRight},
						{Type: tea.KeyRight},
						{Type: tea.KeyDown},
						{Type: tea.KeyDown},
						{Type: tea.KeySpace},
						{Type: tea.KeyDown},
						{Type: tea.KeySpace},
						{Type: tea.KeyEnter},
					},
					ExpectView:            "? Please select your favorite letters:",
					ExpectSelectedIndexes: []int{2, 3, 14, 15},
					ExpectSelectedValues:  []string{"C", "D", "O", "P"},
				},
			},
		},
		{
			name: "single select",
			selectionModel: func() Model {
				m := New()
				m.Prompt = "Choose a color:"
				m.MaxSelections = 1
				m.Choices = []string{"Red", "Green", "Blue"}
				return m
			}(),
			states: []state{
				{
					Name:                  "selects item",
					Inputs:                []tea.KeyMsg{{Type: tea.KeyDown}, {Type: tea.KeySpace}},
					ExpectView:            "? Choose a color:",
					ExpectSelectedIndexes: []int{1},
					ExpectSelectedValues:  []string{"Green"},
				},
			},
		},
	}
	for _, tt := range tests {
		for _, s := range tt.states {
			t.Run(fmt.Sprintf("%s_%s", tt.name, s.Name), func(t *testing.T) {
				tm := teatest.NewTestModel(
					t, &tt.selectionModel,
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

				model := tm.FinalModel(t).(*Model)
				assert.Equal(t, s.ExpectSelectedIndexes, model.SelectedIndexes())
				assert.Equal(t, s.ExpectSelectedValues, model.SelectedValues())

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
