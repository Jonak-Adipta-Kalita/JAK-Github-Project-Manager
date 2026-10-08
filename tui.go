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
	app lipgloss.Style
}

type model struct {
	styles        styles
	height, width int
	list          list.Model
}

func initialModel() model {
	m := model{}
	m.styles = styles{
		app: lipgloss.NewStyle().Padding(1, 2),
	}

	projects := GetProjectsDirs()
	projectItems := make([]list.Item, len(projects))
	for i, project := range projects {
		projectItems[i] = project
	}

	delegate := list.NewDefaultDelegate()

	projectList := list.New(projectItems, delegate, 0, 0)
	projectList.SetShowTitle(false)
	projectList.DisableQuitKeybindings()

	m.list = projectList

	return m
}

func (m *model) updateListProperties() {
	h, v := m.styles.app.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel

	return m, cmd
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
