package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type FileName string

const (
	RENAME_FILE       FileName = "scripts/rename.json"
	GROUP_SORTER_FILE FileName = "scripts/group_sorter.json"
	MAKE_PUBLIC_FILE  FileName = "scripts/make_public.json"
	REPO_REMOVE_FILE  FileName = "scripts/repo_remove.json"
)

type Mode string

const (
	ModeRename      Mode = "RENAME"
	ModeGroupSorter Mode = "GROUP_SORTER"
	ModeMakePublic  Mode = "MAKE_PUBLIC"
	ModeRepoRemove  Mode = "REPO_REMOVE"
)

type Entry interface {
	isArg()
}

type RenameArgs struct {
	OldPath string `json:"old_path"`
	NewName string `json:"new_name"`
}

func (RenameArgs) isArg() {}

type GroupSorterArgs struct {
	Selected  []string `json:"selected"`
	GroupName string   `json:"group_name"`
}

func (GroupSorterArgs) isArg() {}

type MakePublicArgs struct {
	Selected []string `json:"selected"`
}

func (MakePublicArgs) isArg() {}

type RepoRemoveArgs struct {
	Selected []string `json:"selected"`
	DestPath string   `json:"dest_path"`
}

func (RepoRemoveArgs) isArg() {}

func loadScript[T Entry](filePath FileName) ([]T, error) {
	data, err := os.ReadFile(string(filePath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var entries []T
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func appendToScript[T Entry](filePath FileName, arg T) {
	err := os.MkdirAll(string(filePath), 0755)
	if err != nil {
		log.Printf("Error in creating parent folder: %v", err)
		return
	}

	entries, err := loadScript[T](filePath)
	if err != nil {
		log.Printf("Error in loading script: %v", err)
		return
	}

	entries = append(entries, arg)

	out, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		log.Printf("Error marshalling entries: %v", err)
		return
	}

	if err := os.WriteFile(string(filePath), out, 0644); err != nil {
		log.Printf("Error writing to file: %v", err)
		return
	}

	fmt.Println("Stored to Script file")
}

func PutInScript(args Entry) {
	switch a := args.(type) {
	case RepoRemoveArgs:
		appendToScript(REPO_REMOVE_FILE, a)
	case MakePublicArgs:
		appendToScript(MAKE_PUBLIC_FILE, a)
	case GroupSorterArgs:
		appendToScript(GROUP_SORTER_FILE, a)
	case RenameArgs:
		appendToScript(RENAME_FILE, a)
	default:
		log.Fatal("unknown entry type given to PutInScript")
	}
}

func runFile[T Entry](filePath FileName, run func(T)) error {
	entries, err := loadScript[T](filePath)
	if err != nil {
		return err
	}
	for _, e := range entries {
		run(e)
	}
	if err := os.Remove(string(filePath)); err != nil && !os.IsNotExist(err) {
		return err
	}
	fmt.Println("Executed", filePath, "and cleared the file")
	return nil
}

func ExecuteScript() {
	steps := []func() error{
		func() error {
			return runFile(REPO_REMOVE_FILE, func(a RepoRemoveArgs) { GithubRepoRemove(a.Selected, a.DestPath) })
		},
		func() error {
			return runFile(GROUP_SORTER_FILE, func(a GroupSorterArgs) { GroupSorter(a.Selected, a.GroupName) })
		},
		func() error {
			return runFile(MAKE_PUBLIC_FILE, func(a MakePublicArgs) { MakePublic(a.Selected) })
		},
		func() error {
			return runFile(RENAME_FILE, func(a RenameArgs) { RenameProject(a.OldPath, a.NewName) })
		},
	}

	for _, step := range steps {
		if err := step(); err != nil {
			log.Printf("Error executing script: %v", err)
			return
		}
	}
}
