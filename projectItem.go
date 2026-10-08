package main

type projectItem struct {
	title string
	path  string
}

func (i projectItem) Title() string       { return i.title }
func (i projectItem) Description() string { return i.path }
func (i projectItem) FilterValue() string { return i.title }
