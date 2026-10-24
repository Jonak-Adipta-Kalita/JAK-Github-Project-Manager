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

func loadScript(filePath FileName) ([]json.RawMessage, error) {
	data, err := os.ReadFile(string(filePath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var entries []json.RawMessage
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func PutInScript(mode Mode, args any) {
	var filePath FileName

	switch mode {
	case ModeRepoRemove:
		filePath = REPO_REMOVE_FILE
	case ModeMakePublic:
		filePath = MAKE_PUBLIC_FILE
	case ModeGroupSorter:
		filePath = GROUP_SORTER_FILE
	case ModeRename:
		filePath = RENAME_FILE
	default:
		log.Fatal("wrong mode given to put in script")
	}

	entries, err := loadScript(filePath)
	if err != nil {
		log.Printf("Error in loading script: %v", err)
	}

	newEntry, err := json.Marshal(args)
	if err != nil {
		log.Printf("Error in marshalling Args: %v", err)
	}

	entries = append(entries, newEntry)

	out, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		log.Printf("Error marshalling entries: %v", err)
	}

	err = os.WriteFile(string(filePath), out, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}

	fmt.Println("Stored to Script file")
}

func ExecuteScript() {
	// entries, err := loadScript()
	// if err != nil {
	// 	log.Printf("Error in loading script: %v", err)
	// 	return
	// }
	//
	// for _, entry := range entries {
	// 	switch entry.Mode {
	// 	case ModeRepoRemove:
	// 		var a RepoRemoveArgs
	// 		if err := json.Unmarshal(entry.Args, &a); err != nil {
	// 			log.Printf("Error in unmarshaling args: %v", err)
	// 			return
	// 		}
	// 		GithubRepoRemove(a.Selected, a.DestPath)
	// 	case ModeMakePublic:
	// 		var a MakePublicArgs
	// 		if err := json.Unmarshal(entry.Args, &a); err != nil {
	// 			log.Printf("Error in unmarshaling args: %v", err)
	// 			return
	// 		}
	// 		MakePublic(a.Selected)
	// 	case ModeGroupSorter:
	// 		var a GroupSorterArgs
	// 		if err := json.Unmarshal(entry.Args, &a); err != nil {
	// 			log.Printf("Error in unmarshaling args: %v", err)
	// 			return
	// 		}
	// 		GroupSorter(a.Selected, a.GroupName)
	// 	case ModeRename:
	// 		var a RenameArgs
	// 		if err := json.Unmarshal(entry.Args, &a); err != nil {
	// 			log.Printf("Error in unmarshaling args: %v", err)
	// 			return
	// 		}
	// 		renameProject(a.OldPath, a.NewName)
	// 	}
	// }
	//
	// err = os.Remove(filePath)
	// if err != nil && !os.IsNotExist(err) {
	// 	log.Printf("Error in removing script: %v", err)
	// 	return
	// }
	//
	// fmt.Println("Executed the Scripts! And cleared the file for future usage!")
}
