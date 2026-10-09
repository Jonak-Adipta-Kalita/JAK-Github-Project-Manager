package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

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

func ReplaceName(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}

func GetProjectsDirs() (projects []projectItem) {
	path := os.Getenv("PROJECT_DIR")

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range entries {
		projects = append(projects, projectItem{
			title: f.Name(),
			path:  path + f.Name(),
		})
	}

	return projects
}

func main() {
	fmt.Print("Which Mode do you want to enter?\n1) Group Sorter\n2) Remover\n3) Git Clone\n>> ")
	var mode uint8
	fmt.Scan(&mode)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if mode == 3 {
		CloneRepositories()
	} else {
		selected := Tui()
		fmt.Println()

		if len(selected) == 0 {
			fmt.Println("Select the Projects too...")
			os.Exit(1)
		} else {
			fmt.Printf("Selected: %v", selected)
		}

		if mode == 1 {
		} else if mode == 2 {
			GithubRemover(selected, filepath.Join(os.Getenv("REMOVED_DIR"), "Stolen/"))
		}
	}
}
