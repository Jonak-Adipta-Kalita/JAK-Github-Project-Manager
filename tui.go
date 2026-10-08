package main

import (
	// "fmt"
	"log"

	// "charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
}

func newStyles() styles {
	return styles{
		app: lipgloss.NewStyle().
			Padding(1, 2),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		statusMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")),
	}
}

type model struct {
	styles        styles
	height, width int
	list          list.Model
}

func initialModel() model {
	m := model{}
	m.styles = newStyles()

	projects := GetProjectsDirs()
	projectItems := make([]list.Item, len(projects))
	for i, project := range projects {
		projectItems[i] = project
	}

	delegate := list.NewDefaultDelegate()

	projectList := list.New(projectItems, delegate, 0, 0)
	projectList.Title = "Choose Projects"
	projectList.Styles.Title = m.styles.title

	m.list = projectList

	return m
}

func (m *model) updateListProperties() {
	h, v := m.styles.app.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)
	m.list.Styles.Title = m.styles.title
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	v := tea.NewView(m.styles.app.Render(m.list.View()))
	v.AltScreen = true
	return v
}

func Tui() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error in TUI: %v", err)
	}
}
