package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func checkOfflineMode() bool {
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 3*time.Second)
	if err != nil {
		log.Printf("error in hitting cloudflare: %v", err)
		return true
	}

	defer conn.Close()
	return false
}

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

func GroupSorter(projectPaths []string, groupName string) {
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
		log.Printf("Failed to create repository %s\n", repo)
	}
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

func main() {
	var scriptMode bool

	if checkOfflineMode() {
		scriptMode = true
		fmt.Println("In Offline Mode, so changes will be saved to a script!\n")
	} else {
		fmt.Print("1) Script Mode or 2) Interactive Mode? >> ")
		var offlineInput uint8
		fmt.Scan(&offlineInput)
		fmt.Print("\n")

		if offlineInput == 2 {
			scriptMode = false
			fmt.Println("Using Interactive Mode")
		} else {
			scriptMode = true
			fmt.Println("Using Script Mode")
		}
	}
	fmt.Println()

	fmt.Print("Which Mode do you want to enter?\n\n1) Git Clone\n2) Just Browse\n\n3) Group Sorter \n4) Make Public\n5) Remove\n\n6) Run Offline Scripts\n\n>> ")
	var mode uint8
	fmt.Scan(&mode)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if !slices.Contains([]uint8{1, 2, 3, 4, 5}, mode) {
		log.Fatal("Select one of 1, 2, 3, 4, 5")
	}

	if mode == 1 {
		if scriptMode {
			fmt.Println("Come back later with Internet!")
			return
		}

		CloneRepositories()
	} else {
		selected := Tui(scriptMode)
		fmt.Println()

		if mode != 2 && len(selected) == 0 {
			fmt.Println("Select the Projects too...")
			os.Exit(1)
		} else {
			fmt.Printf("Selected: %v\n\n", selected)
		}

		if mode == 3 {
			var groupName string
			fmt.Print("Group Name (No Whitespace): ")
			fmt.Scan(&groupName)

			if scriptMode {
				PutInOfflineScript(ModeGroupSorter, GroupSorterArgs{
					Selected:  selected,
					GroupName: groupName,
				})
				return
			}

			GroupSorter(selected, groupName)
		} else if mode == 4 {
			if scriptMode {
				PutInOfflineScript(ModeMakePublic, MakePublicArgs{Selected: selected})
				return
			}

			MakePublic(selected)
		} else if mode == 5 {
			if scriptMode {
				PutInOfflineScript(ModeRepoRemove, RepoRemoveArgs{Selected: selected, DestPath: os.Getenv("REMOVED_DIR")})
				return
			}

			GithubRepoRemove(selected, os.Getenv("REMOVED_DIR"))
		}
	}
}
