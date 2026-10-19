package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Mode string

const (
	ModeRename      Mode = "RENAME"
	ModeGroupSorter Mode = "GROUP_SORTER"
	ModeMakePublic  Mode = "MAKE_PUBLIC"
	ModeRepoRemove  Mode = "REPO_REMOVE"
)

func ExecuteScript() {
	f, err := os.Open("offline-mode.jak")
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		fmt.Printf("%d: %s\n", lineNum, line)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file at line %d: %v", lineNum+1, err)
	}
}

func PutInScript(mode Mode, args ...any) {
	content := fmt.Sprintf("%s", mode)

	for _, a := range args {
		content += fmt.Sprintf(" %v", a)
	}
	content += "\n"

	f, err := os.OpenFile("offline-mode.jak", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		log.Printf("Error writing to file: %v", err)
		return
	}

	fmt.Println("Stored to Offline-Mode file")
}
