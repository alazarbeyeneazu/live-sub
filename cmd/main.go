package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hacker301et/live-sub/view"
)

func main() {
	m := view.NewView()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}

}
