package main

import (
	"errors"
	"regexp"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

var validName = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func validateName(name, old string) error {
	switch {
	case name == "":
		return errors.New("name can't be empty")
	case name == old:
		return errors.New("name is unchanged")
	case !validName.MatchString(name):
		return errors.New("use only letters, digits, - _ .")
	}
	return nil
}

type renameModalTUI struct {
	model   textinput.Model
	oldName string
	newName string
	error   string
}

func (tui renameModalTUI) Setup() {
	tui.model.SetWidth(40)
	tui.model.CharLimit = 100
	// tui.model.Placeholder = //TODO
}

func (tui renameModalTUI) UpdateRename(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.renaming = false
			tui.error = ""
			tui.model.Blur()
			return m, nil
		}
	}

	var cmd tea.Cmd
	tui.model, cmd = tui.model.Update(msg)

	return m, cmd
}
