package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	app   lipgloss.Style
	modal lipgloss.Style
	title lipgloss.Style
}

type keyMap struct {
	toggle  key.Binding
	open    key.Binding
	confirm key.Binding
	rename  key.Binding
}

type model struct {
	styles        styles
	keys          keyMap
	height, width int
	list          list.Model
	renameModal   renameModalTUI
	selected      map[string]struct{}
	confirmed     bool
	offlineMode   bool
}

func initialModel(offlineMode bool) model {
	m := model{}
	m.offlineMode = offlineMode

	m.styles = styles{
		app: lipgloss.NewStyle().Padding(1, 2),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		modal: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#25A065")).
			Padding(1, 2),
	}

	m.keys = keyMap{
		toggle: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "toggle"),
		),
		open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open in yazi"),
		),
		confirm: key.NewBinding(
			key.WithKeys("alt+o"),
			key.WithHelp("alt+o", "confirm"),
		),
		rename: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename in github"),
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
	projectList.Title = "Projects"
	projectList.Styles.Title = m.styles.title
	projectList.StatusMessageLifetime = 4 * time.Second

	projectList.DisableQuitKeybindings()

	projectList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{m.keys.toggle, m.keys.open, m.keys.rename, m.keys.confirm}
	}

	m.list = projectList

	renameModal := textinput.New()
	m.renameModal.model = renameModal
	m.renameModal.Setup()

	return m
}

func (m *model) updateListProperties() {
	h, v := m.styles.app.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)
}

func (m model) Init() tea.Cmd {
	return nil
}

func openInYazi(path string) error {
	cmd := exec.Command("kitty", "yazi", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if d, ok := msg.(renameDoneMsg); ok {
		return m.handleRenameDone(d)
	}
	if m.renameModal.renaming {
		return m.updateRename(msg)
	}

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

		case key.Matches(msg, m.keys.open):
			it, ok := m.list.SelectedItem().(projectItem)
			if !ok {
				return m, nil
			}
			if err := openInYazi(it.path); err != nil {
				return m, m.list.NewStatusMessage("Error: " + err.Error())
			}
			cmd := exec.Command("xdg-open", fmt.Sprintf("https://github.com/%s/%s", os.Getenv("GITHUB_USERNAME"), filepath.Base(it.path)))
			if err := cmd.Run(); err != nil {
				log.Printf("Can't open in Browser")
				return m, m.list.NewStatusMessage("Error: " + err.Error())
			}
			return m, m.list.NewStatusMessage("Opened " + it.title + " in yazi & browser")

		case key.Matches(msg, m.keys.rename):
			it, ok := m.list.SelectedItem().(projectItem)
			if !ok {
				return m, nil
			}

			m.renameModal.renaming = true
			m.renameModal.oldName = it.title
			m.renameModal.oldPath = it.path
			m.renameModal.error = ""
			m.renameModal.model.SetValue(it.title)
			m.renameModal.model.CursorEnd()
			focusCmd := m.renameModal.model.Focus()

			return m, focusCmd

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
	var content string
	if m.renameModal.renaming {
		content = m.renameView()
	} else {
		content = m.styles.app.Render(m.list.View())
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func Tui(offlineMode bool) []string {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("can't open log file: %v", err)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel(offlineMode))
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
