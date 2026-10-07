package main

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func initialModel(choices []string) model {
	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", "space":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	s := "Select the Projects\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit\n"

	v := tea.NewView(s)
	v.WindowTitle = "JAK Github Project Manager"
	v.AltScreen = true

	return v
}

func Tui() {
	projects := GetProjectsDirs()

	p := tea.NewProgram(initialModel(projects))
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error in TUI: %v", err)
	}
}
