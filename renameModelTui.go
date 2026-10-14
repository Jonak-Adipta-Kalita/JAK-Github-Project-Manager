package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	renaming bool
	model    textinput.Model
	oldName  string
	oldPath  string
	error    string
	busy     bool
}

func (tui *renameModalTUI) Setup() {
	tui.model.SetWidth(40)
	tui.model.CharLimit = 150
	tui.model.Placeholder = "Meow"
}

func (m model) updateRename(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil

	case tea.KeyPressMsg:
		if m.renameModal.busy {
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.renameModal.renaming = false
			m.renameModal.error = ""
			m.renameModal.model.Blur()
			return m, nil

		case "enter":
			name := strings.TrimSpace(m.renameModal.model.Value())
			if err := validateName(name, m.renameModal.oldName); err != nil {
				m.renameModal.error = err.Error()
				return m, nil
			}

			m.renameModal.busy = true
			m.renameModal.error = ""
			m.renameModal.model.Blur()

			return m, renameCmd(m.renameModal.oldPath, name)
		}
	}

	var cmd tea.Cmd
	m.renameModal.model, cmd = m.renameModal.model.Update(msg)

	return m, cmd
}

func (m model) renameView() string {
	r := m.renameModal

	title := lipgloss.NewStyle().Bold(true).Render("Rename " + r.oldName)
	hint := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("enter confirm • esc cancel")

	lines := []string{title, "", r.model.View()}
	if r.error != "" {
		lines = append(lines, "",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Render(r.error))
	}

	if r.busy {
		lines = append(lines, "", lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render("Renaming..."))
	} else {
		lines = append(lines, "", hint)
	}

	box := m.styles.modal.Render(strings.Join(lines, "\n"))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

type renameDoneMsg struct {
	oldPath string
	newPath string
	newName string
	err     error
}

func renameProject(oldPath, newName string) (string, error) {
	newPath := filepath.Join(filepath.Dir(oldPath), newName)
	if _, err := os.Stat(newPath); err == nil {
		return "", fmt.Errorf("folder %q already exists", newName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", "repo", "rename", "--yes", newName)
	cmd.Dir = oldPath

	if out, err := cmd.CombinedOutput(); err != nil {
		msg := strings.Join(strings.Fields(string(out)), " ")
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New("gh: " + msg)
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return "", fmt.Errorf("renamed on GitHub, but folder rename failed: %w", err)
	}

	return newPath, nil
}

func renameCmd(oldPath, newName string) tea.Cmd {
	return func() tea.Msg {
		newPath, err := renameProject(oldPath, newName)

		return renameDoneMsg{
			oldPath: oldPath,
			newPath: newPath,
			newName: newName,
			err:     err,
		}
	}
}

func (m model) handleRenameDone(d renameDoneMsg) (tea.Model, tea.Cmd) {
	m.renameModal.busy = false

	if d.err != nil {
		m.renameModal.error = d.err.Error()
		focusCmd := m.renameModal.model.Focus()
		return m, focusCmd
	}

	m.renameModal.renaming = false
	m.renameModal.error = ""

	if _, wasSelected := m.selected[d.oldPath]; wasSelected {
		delete(m.selected, d.oldPath)
		m.selected[d.newPath] = struct{}{}
	}

	var setCmd tea.Cmd
	for i, it := range m.list.Items() {
		p, ok := it.(projectItem)
		if !ok || p.path != d.oldPath {
			continue
		}
		p.title = d.newName
		p.path = d.newPath
		setCmd = m.list.SetItem(i, p)
		break
	}

	statusCmd := m.list.NewStatusMessage("Renamed to " + d.newName)
	return m, tea.Batch(setCmd, statusCmd)
}
