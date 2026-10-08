package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func ReplaceName(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}

func GetProjectsDirs() (projects []projectItem) {
	path := "/home/jonak/Projects/Sorting/Main-Clones/"

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

func main() {
	fmt.Print("Which Mode do you want to enter?\n1) Group Sorter\n2) Remover\n3) Git Clone\n>> ")
	var mode uint
	fmt.Scan(&mode)

	if mode == 3 {
		CloneRepositories()
	} else {
		Tui()
	}
}
