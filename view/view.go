package view

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hacker301et/live-sub/internal"
	"github.com/hacker301et/live-sub/models"
)

type endLoaidng struct{}

const (
	SUBLISTER = "sub-lister"
	AMASS     = "Amass"
)

type model struct {
	Sub      chan models.ResponseMsg
	spinner  spinner.Model
	webtable table.Model
	apiTable table.Model
	FQDN     textinput.Model
	typing   bool
	err      error
	rows     []table.Row
	apiRows  []table.Row
}

func NewView() *model {
	sp := spinner.New()
	sp.Tick()
	sp.Spinner = spinner.Points
	ti := textinput.NewModel()
	ti.Focus()
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Border(lipgloss.BlockBorder())

	ti.Placeholder = "Enter FQDM here "
	m := model{
		Sub:     make(chan models.ResponseMsg),
		FQDN:    ti,
		typing:  true,
		rows:    make([]table.Row, 0),
		spinner: sp,
		apiRows: make([]table.Row, 0),
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

	t2 := table.New(
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
	t2.SetStyles(s)
	m.webtable = t
	m.apiTable = t2
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		subs := internal.SubLister(m.FQDN.Value())
		if len(subs) > 0 {
			internal.CheckSubDomain(subs, m.Sub, SUBLISTER)
		}

	}(&wg)
	amassSubs := internal.AmassFindSubDomains(m.FQDN.Value())
	m.endLoading()
	internal.CheckSubDomain(amassSubs, m.Sub, AMASS)
	wg.Wait()

}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch teaMsg := msg.(type) {
	case tea.KeyMsg:
		switch teaMsg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.typing {
				m.typing = false
				go m.FindSubDomains()
				return m, tea.Batch(cmd, m.spinner.Tick)
			}
		case "up":
			m.webtable, cmd = m.webtable.Update(msg)
			return m, cmd
		case "down":
			m.webtable, cmd = m.webtable.Update(msg)
			return m, cmd
		case "left":
			m.apiTable, cmd = m.apiTable.Update(tea.KeyMsg{Type: tea.KeyUp})
			return m, cmd
		case "right":
			m.apiTable, cmd = m.apiTable.Update(tea.KeyMsg{Type: tea.KeyDown})
			return m, cmd

		}
	case models.ResponseMsg:
		if strings.Contains(teaMsg.FQDN, "api") {
			m.apiRows = append(m.apiRows, table.Row{teaMsg.ToolName, teaMsg.FQDN})
			m.apiTable.SetRows(m.apiRows)
			return m, waitForActivity(m.Sub)
		}
		m.rows = append(m.rows, table.Row{teaMsg.ToolName, teaMsg.FQDN})
		m.webtable.SetRows(m.rows)
		return m, waitForActivity(m.Sub)
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.webtable.SetWidth(teaMsg.Width)
		m.webtable.SetHeight(teaMsg.Height/2 - 10)
	}

	m.FQDN, cmd = m.FQDN.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.typing {
		return m.FQDN.View()
	}
	// if m.err != nil {
	// 	tea.Quit()
	// 	return "Error while getting the host"
	// }
	return baseStyle.Render("      "+m.spinner.View()+"\n"+m.webtable.View()) + "\n" + baseStyle.UnsetAlignHorizontal().Render(" APIs     "+"\n"+m.apiTable.View()) + "\n"
}
