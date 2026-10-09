package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

func getClonedDirs() []string {
	projects := GetProjectsDirs()
	cloned := make([]string, len(projects))

	for _, project := range projects {
		cloned = append(cloned, project.Title())
	}

	return cloned
}

func CloneRepositories() {
	original, err := os.ReadDir(os.Getenv("LIST_DIR"))
	if err != nil {
		log.Fatal(err)
	}

	cloned := getClonedDirs()

	err = os.MkdirAll(filepath.Dir(os.Getenv("PROJECT_DIR")), 0755)
	if err != nil {
		fmt.Printf("Failed to create destination directories: %v\n", err)
		return
	}

	var failed []string

	for _, f := range original {
		var name string

		if f.IsDir() {
			name = ReplaceName(f.Name())
		} else {
			name = ReplaceName(strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())))
		}

		if slices.Contains(cloned, name) {
			continue
		}

		fmt.Printf("Cloning %s\n", name)

		url := fmt.Sprintf("git@github.com:%s/%s.git", os.Getenv("GITHUB_USERNAME"), name)
		cmd := exec.Command("git", "-C", os.Getenv("PROJECT_DIR"), "clone", url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Failed to clone %s: %v\n\n", name, err)

			failed = append(failed, name)
		} else {
			fmt.Printf("Cloned %s\n\n", name)

			if f.IsDir() {
				err := os.RemoveAll(fmt.Sprintf("%s/%v", os.Getenv("LIST_DIR"), f.Name()))
				if err != nil {
					fmt.Printf("Failed to delete folder %s: %v\n", f.Name(), err)
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("Failed: %v", failed)
}
