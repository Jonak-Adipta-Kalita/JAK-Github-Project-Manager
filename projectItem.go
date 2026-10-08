package main

type projectItem struct {
	title    string
	path     string
	selected map[string]struct{}
}

func (i projectItem) Title() string {
	if _, ok := i.selected[i.path]; ok {
		return "[x] " + i.title
	}
	return "[ ] " + i.title
}

func (i projectItem) Description() string { return i.path }
func (i projectItem) FilterValue() string { return i.title }
