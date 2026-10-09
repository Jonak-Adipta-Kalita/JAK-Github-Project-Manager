package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func GithubRemover(projectPaths []string, folderPath string) {
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

		cmd := exec.Command("gh", "repo", "delete", repo, "--yes")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Failed to delete repository %s", repo)
		}
	}

	fmt.Println()
	fmt.Printf("Failed: %v\n", failed)
}
