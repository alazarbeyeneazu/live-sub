package view

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hacker301et/live-sub/internal"
	"github.com/hacker301et/live-sub/models"
)

// type models.ResponseMsg struct{}
type endLoaidng struct{}
type model struct {
	Sub    chan models.ResponseMsg
	table  table.Model
	FQDN   textinput.Model
	typing bool
	rows   []table.Row
}

func NewView() *model {

	ti := textinput.NewModel()
	ti.Focus()
	ti.Placeholder = "Enter FQDM here "
	m := model{
		Sub:    make(chan models.ResponseMsg),
		FQDN:   ti,
		typing: true,
		rows:   make([]table.Row, 0),
	}
	columns := []table.Column{
		{Title: "Tool", Width: 10},
		{Title: "Live Sub", Width: 150},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(m.rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	m.table = t
	return &m
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func waitForActivity(sub chan models.ResponseMsg) tea.Cmd {

	return func() tea.Msg {
		return models.ResponseMsg(<-sub)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.Sub),
	)
}
func (m *model) endLoading() tea.Cmd {

	return func() tea.Msg {
		return endLoaidng{}
	}
}
func (m *model) FindSubDomains() {
	subs := internal.SubLister(m.FQDN.Value())
	m.endLoading()
	internal.CheckSubDomain(subs, m.Sub)
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch teaMsg := msg.(type) {
	case tea.KeyMsg:
		switch teaMsg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.typing = false
			go m.FindSubDomains()
			return m, cmd
		case "up":
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		case "down":
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		}
	case models.ResponseMsg:
		m.rows = append(m.rows, table.Row{teaMsg.ToolName, teaMsg.FQDN})
		m.table.SetRows(m.rows)
		return m, waitForActivity(m.Sub)
	}
	m.FQDN, cmd = m.FQDN.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.typing {
		return m.FQDN.View()
	}
	return baseStyle.Render(m.table.View()) + "\n"
}
