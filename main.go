package main

import (
	"fmt"
	"log"
	"os"
	"strings"
);

func ReplaceName(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}

func main() {
	entries, err := os.ReadDir("/home/jonak/Projects/Sorting/Main-Clones/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Which Mode do you want to enter?\n1) Group Sorter\n2) Remover\n3) Git Clone\n>> ")
	var mode uint
	fmt.Scan(&mode)

	if mode == 3 {
		CloneRepositories(entries)
	} else {

	}
}

