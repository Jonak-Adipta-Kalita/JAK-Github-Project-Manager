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

func GithubRepoRemove(projectPaths []string, destPath string) {
	var failed []string

	err := os.MkdirAll(destPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create destination directories: %v\n", err)
		return
	}

	for _, projectPath := range projectPaths {
		name := filepath.Base(projectPath)

		err := os.Rename(projectPath, filepath.Join(destPath, name))
		if err != nil {
			fmt.Printf("Error moving folder: %v\n", err)
			failed = append(failed, projectPath)
			continue
		}

		repo := fmt.Sprintf("%s/%s", os.Getenv("GITHUB_USERNAME"), name)

		cmd := exec.Command("gh", "repo", "delete", repo, "--yes")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Failed to delete repository %s", repo)
			failed = append(failed, projectPath)
		}
	}

	fmt.Println()
	fmt.Printf("Failed: %v\n", failed)
}

func GroupSorter(projectPaths []string) {
	var groupName string
	fmt.Print("Group Name (No Whitespace): ")
	fmt.Scan(&groupName)

	groupDir := filepath.Join(os.Getenv("GROUP_DIR"), groupName)

	GithubRepoRemove(projectPaths, groupDir)

	for _, p := range projectPaths {
		os.RemoveAll(filepath.Join(groupDir, filepath.Base(p), ".git"))
	}

	git := func(args ...string) error {
		c := exec.Command("git", args...)
		c.Dir = groupDir
		c.Stdout, c.Stderr = os.Stdout, os.Stderr
		return c.Run()
	}

	for i, args := range [][]string{
		{"init", "-b", "main"},
		{"add", "."},
		{"commit", "-m", "Grouped Projects <3"},
		{"push", "origin", "main"},
	} {
		if dirExists(groupDir) && i == 0 {
			continue
		} else if i == 3 {
			continue
		}

		if err := git(args...); err != nil {
			log.Printf("git %s failed: %v", args[0], err)
			return
		}
	}

	if dirExists(groupDir) {
		return
	}

	repo := fmt.Sprintf("%s/%s", os.Getenv("GITHUB_USERNAME"), groupName)
	cmd := exec.Command("gh", "repo", "create", repo, "--private", "--source=.", "--remote=origin", "--push")

	cmd.Dir = groupDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Failed to create repository %s", repo)
	}
}

func main() {
	fmt.Print("Which Mode do you want to enter?\n1) Group Sorter\n2) Remover\n3) Git Clone\n4) Just Browse Projects\n>> ")
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
			fmt.Printf("Selected: %v\n\n", selected)
		}

		if mode == 1 {
			GroupSorter(selected)
		} else if mode == 2 {
			GithubRepoRemove(selected, filepath.Join(os.Getenv("REMOVED_DIR"), "Stolen/"))
		} else if mode == 4 {
			Tui()
		}
	}
}
