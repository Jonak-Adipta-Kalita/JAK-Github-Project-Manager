package main

import (
	"fmt"
	"log"
	"sort"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	app lipgloss.Style
}

type keyMap struct {
	toggle  key.Binding
	confirm key.Binding
}

type model struct {
	styles        styles
	keys          keyMap
	height, width int
	list          list.Model
	selected      map[string]struct{}
	confirmed     bool
}

func initialModel() model {
	m := model{}
	m.styles = styles{
		app: lipgloss.NewStyle().Padding(1, 2),
	}

	m.keys = keyMap{
		toggle: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "toggle"),
		),
		confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
	}

	m.selected = make(map[string]struct{})

	projects := GetProjectsDirs()
	projectItems := make([]list.Item, len(projects))
	for i, project := range projects {
		project.selected = m.selected
		projectItems[i] = project
	}

	delegate := list.NewDefaultDelegate()

	projectList := list.New(projectItems, delegate, 0, 0)
	projectList.SetShowTitle(false)
	projectList.DisableQuitKeybindings()
	projectList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{m.keys.toggle, m.keys.confirm}
	}

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
		switch {
		case key.Matches(msg, m.keys.toggle):
			it, ok := m.list.SelectedItem().(projectItem)
			if !ok {
				return m, nil
			}
			if _, picked := m.selected[it.path]; picked {
				delete(m.selected, it.path)
			} else {
				m.selected[it.path] = struct{}{}
			}
			m.list.CursorDown()
			status := m.list.NewStatusMessage(fmt.Sprintf("%d selected", len(m.selected)))
			return m, status
		case key.Matches(msg, m.keys.confirm):
			m.confirmed = true
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

func Tui() []string {
	p := tea.NewProgram(initialModel())
	final, err := p.Run()
	if err != nil {
		log.Fatalf("Error in TUI: %v", err)
	}

	m, ok := final.(model)
	if !ok || !m.confirmed {
		return nil
	}

	paths := make([]string, 0, len(m.selected))
	for p := range m.selected {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}
