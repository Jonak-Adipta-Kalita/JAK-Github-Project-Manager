package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
)

func getClonedDirs(entries []os.DirEntry) (projects []string) {
	for _, f := range entries {
		projects = append(projects, f.Name())
	}

	return projects
}

func CloneRepositories(entries []os.DirEntry) {
	cloned := getClonedDirs(entries)

	original, err := os.ReadDir("/home/jonak/Projects/Sorting/Main")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range original {
		if !f.IsDir() {
			continue
		}

		name := ReplaceName(f.Name())
		if slices.Contains(cloned, name) {
			continue
		}

		fmt.Printf("Cloning %s\n", name)

		url := fmt.Sprintf("git@github.com:Jonak-Adipta-Kalita/%s.git", name)
		cmd := exec.Command("git", "-C", "/home/jonak/Projects/Sorting/Main-Clones", "clone", url)

		if err := cmd.Run(); err != nil {
			log.Printf("Failed to clone %s: %v\n", name, err)
		} else {
			fmt.Printf("Cloned %s\n\n", name)
			err := os.RemoveAll(fmt.Sprintf("/home/jonak/Projects/Sorting/Main/%v", f.Name()))
			if err != nil {
				fmt.Printf("Failed to delete folder %s: %v\n", f.Name(), err)
			}
		}
	}
}
