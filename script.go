package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var filePath string = "offline-mode.json"

type Mode string

const (
	ModeRename      Mode = "RENAME"
	ModeGroupSorter Mode = "GROUP_SORTER"
	ModeMakePublic  Mode = "MAKE_PUBLIC"
	ModeRepoRemove  Mode = "REPO_REMOVE"
)

type Entry struct {
	Mode Mode            `json:"mode"`
	Args json.RawMessage `json:"args"`
}

type RenameArgs struct {
	OldPath string `json:"old_path"`
	NewName string `json:"new_name"`
}

type GroupSorterArgs struct {
	Selected  []string `json:"selected"`
	GroupName string   `json:"group_name"`
}

type MakePublicArgs struct {
	Selected []string `json:"selected"`
}

type RepoRemoveArgs struct {
	Selected []string `json:"selected"`
	DestPath string   `json:"dest_path"`
}

func loadOfflineScript() ([]Entry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func ExecuteOflineScript() {

}

func PutInOfflineScript(mode Mode, args any) {
	entries, err := loadOfflineScript()
	if err != nil {
		log.Printf("Error in loading offling-mode script: %v", err)
	}

	marshalledArgs, err := json.Marshal(args)
	if err != nil {
		log.Printf("Error in marshalling Args: %v", err)
	}

	newEntry := Entry{
		Mode: mode,
		Args: marshalledArgs,
	}

	entries = append(entries, newEntry)

	out, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		log.Printf("Error marshalling entries: %v", err)
	}

	err = os.WriteFile(filePath, out, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}

	fmt.Println("Stored to Offline-Mode file")
}
