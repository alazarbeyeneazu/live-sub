package view

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	csv "github.com/gocarina/gocsv"
	"github.com/hacker301et/live-sub/internal"
	"github.com/hacker301et/live-sub/models"
)

type endLoaidng struct{}

const (
	SUBLISTER = "sub-lister"
	AMASS     = "Amass"
	SUBFINDER = "Subfinder"
	STARTDATE = 30
)

type model struct {
	Sub        chan models.ResponseMsg
	spinner    spinner.Model
	webtable   table.Model
	apiTable   table.Model
	FQDN       textinput.Model
	typing     bool
	err        error
	rows       []table.Row
	apiRows    []table.Row
	apiFound   bool
	rowTracker map[string]bool
}

func (m *model) readCSV() ([]models.ResponseMsg, error) {
	fileName := fmt.Sprintf("scans/%s", m.FQDN.Value())
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return []models.ResponseMsg{}, err
	}
	defer file.Close()
	var domains []models.ResponseMsg
	if err := csv.UnmarshalFile(file, &domains); err != nil {
		return []models.ResponseMsg{}, err
	}
	return domains, nil

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
		Sub:        make(chan models.ResponseMsg),
		FQDN:       ti,
		typing:     true,
		rows:       make([]table.Row, 0),
		spinner:    sp,
		apiRows:    make([]table.Row, 0),
		rowTracker: make(map[string]bool),
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
		BorderForeground(lipgloss.Color("34")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("234")).
		Background(lipgloss.Color("34")).
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
	go m.endLoading()
	var wg sync.WaitGroup
	wg.Add(3)
	go func(w *sync.WaitGroup) {
		defer w.Done()

		for date := STARTDATE; date > -1; date-- {
			subs := internal.SubLister(m.FQDN.Value(), date)
			if len(subs) > 0 {
				internal.CheckSubDomain(subs, m.Sub, SUBLISTER)
				break
			}
		}

	}(&wg)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		amassSubs := internal.AmassFindSubDomains(m.FQDN.Value())
		internal.CheckSubDomain(amassSubs, m.Sub, AMASS)
	}(&wg)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		subfinderSubs := internal.SubFinderFindSubDomains(m.FQDN.Value())
		internal.CheckSubDomain(subfinderSubs, m.Sub, SUBFINDER)
	}(&wg)
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
		if exist := m.rowTracker[teaMsg.FQDN]; exist {
			return m, nil
		}
		m.rowTracker[teaMsg.FQDN] = true
		if strings.Contains(teaMsg.FQDN, "api") {
			m.apiFound = true
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
		m.webtable.SetHeight(teaMsg.Height/2 - 10)
		m.apiTable.SetHeight(teaMsg.Height/2 - 10)
	}

	m.FQDN, cmd = m.FQDN.Update(msg)

	return m, cmd
}
func returnView(views ...string) string {
	var v string
	for _, view := range views {
		v = v + view
	}
	return v
}
func (m model) View() string {
	webView := baseStyle.Render("      "+m.spinner.View()+"\n Wesbites  ⬆  To Move Up   ⬇  To Move down \n"+m.webtable.View()) + "\n"
	apiView := baseStyle.Render(" APIs     "+"  ⬅  To Move Up   ➡  To Move down \n"+m.apiTable.View()) + "\n"
	if m.typing {
		return m.FQDN.View()
	}

	return returnView(webView, apiView)

}
