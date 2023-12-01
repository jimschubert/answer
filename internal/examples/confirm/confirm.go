package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jimschubert/answer/confirm"
)

var vertical = flag.Bool("vertical", false, "use vertical selection")
var selectable = flag.Bool("selectable", false, "selectable items")

func main() {
	flag.Parse()

	m := confirm.New()
	m.Prompt = "Do you like pie?"

	if *selectable {
		m.AcceptedDecisionText = "Yes"
		m.DeniedDecisionText = "No"
		if *vertical {
			m.Rendering = confirm.VerticalSelection
		} else {
			m.Rendering = confirm.HorizontalSelection
		}
	}

	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "You chose (%v)\nOther decision formatting:\n\tm.Selected().YesNoString(): %v\n\tm.Selected().TrueFalseString(): %v\n\tm.Selected().String(): %v\n", m.Value(), m.Selected().YesNoString(), m.Selected().TrueFalseString(), m.Selected().String())

}
