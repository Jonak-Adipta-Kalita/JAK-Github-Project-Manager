package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

func GithubRepoAction(projectPaths []string, folderPath, action string) {
	var failed []string

	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create destination directories: %v\n", err)
		return
	}

	for _, projectPath := range projectPaths {
		name := filepath.Base(projectPath)

		err := os.Rename(projectPath, filepath.Join(folderPath, name))
		if err != nil {
			fmt.Printf("Error moving folder: %v\n", err)
			failed = append(failed, projectPath)
			continue
		}

		repo := fmt.Sprintf("%s/%s", os.Getenv("GITHUB_USERNAME"), name)

		cmd := exec.Command("gh", "repo", action, repo, "--yes")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Failed to %s repository %s", action, repo)
		}
	}

	fmt.Println()
	fmt.Printf("Failed: %v\n", failed)
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
			GroupSorter(selected)
		} else if mode == 2 {
			GithubRepoAction(selected, filepath.Join(os.Getenv("REMOVED_DIR"), "Stolen/"), "delete")
		}
	}
}
