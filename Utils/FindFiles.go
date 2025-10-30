package utils

import (
	"fmt"
	"os"
)

func FindFilesInCurrentDirectory() (fileList []os.FileInfo) {

	// Look for all the file in the current directory, but not the sub folders!
	var files, issue = os.Open(".")
	if issue != nil {
		fmt.Printf("Error opening directory %s\n", issue)
		os.Exit(1)
	}

	// List the files in the folder
	fileList, err := files.Readdir(0)
	if err != nil {
		fmt.Printf("Error reading directory, %s\n", err)
		os.Exit(1)
	}

	return fileList
}

func MakeDirectoryList(fileList []os.FileInfo) []string {

	// Initialise a list of the directories
	directoryList := []string{}

	// Loop through all the files and check if their a directory or not
	for _, i := range fileList {
		if i.IsDir() {
			directoryList = append(directoryList, i.Name())
		}
	}

	return directoryList
}
