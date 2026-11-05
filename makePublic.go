package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func safeSliceGet(slice []string, i int) string {
	if i >= 0 && i < len(slice) {
		return slice[i]
	}

	return ""
}

func MakePublic(projectPaths []string) {
	var failed []string

	for _, projectPath := range projectPaths {
		name := filepath.Base(projectPath)
		repo := fmt.Sprintf("%s/%s", os.Getenv("GITHUB_USERNAME"), name)

		cmd := exec.Command("gh", "repo", "edit", repo, "--visibility", "public", "--accept-visibility-change-consequences")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Failed to make the repository %s public\n", repo)
		}
	}

	fmt.Println()
	fmt.Printf("Failed: %v\n", failed)
}

func JAKToTutorial(projectPaths []string, descriptions []string) {
	var failed []string

	for i, projectPath := range projectPaths {
		name := filepath.Base(projectPath)
		repo := fmt.Sprintf("%s/%s", os.Getenv("GITHUB_USERNAME"), name)

		cmd := exec.Command("gh", "repo", "edit", repo, "--visibility", "public", "--accept-visibility-change-consequences", "--description", safeSliceGet(descriptions, i))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to make the repository %s public\n", repo)
		}

		splitName := strings.Split(name, "-")
		splitName[0] = "Tutorial"
		newName := strings.Join(splitName, "-")
		_, err := RenameProject(projectPath, newName)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println()
	fmt.Printf("Failed: %v\n", failed)
}
