package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {

	if len(os.Args[1:]) >= 1 {
		er := CLI(os.Args[1:])
		if er != nil {
			fmt.Printf("Error parsing the command line argument, %v\n", er)
		} else {
			return
		}
	}

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

	// CHECK to see if their is a git folder
	// Initilaise the known files to ignore!
	unwantedFiles := []string{".localized", ".DS_Store", ".gitignore"}
	unwantedExtentions := []string{".app", ".exe", ".elf", ".md"}

	// Initialise a list of the directories
	directoryList := []string{}

	for _, i := range fileList {
		if i.IsDir() {
			directoryList = append(directoryList, i.Name())
		}
	}

	// Look in the directories for a git folder
	if slices.Contains(directoryList, "git") {
		fmt.Printf("[ERROR]: No git folder found\n")
		os.Exit(1)
	}

	// Check there is an origin, and exit if not!
	_, remoteOriginErr := GetRemoteOrigin()
	if remoteOriginErr != nil {
		fmt.Printf("[ERROR]: %s\n", remoteOriginErr)
		os.Exit(1)
	}

	// Get a list of all current issues
	listOfGithubIssues, githubErr := ListGithubIssues()
	if githubErr != nil {
		if errors.Is(githubErr, fmt.Errorf("there were no github issues")) {
			fmt.Printf("[ERROR]: There was an error getting issues: %v\n", githubErr)
			return
		}
	}

	// Get the number of existing issues
	CurrentNumberOfIssues := len(listOfGithubIssues)

	var foundNewTODO bool = false
	for _, fileName := range fileList {

		if fileName.IsDir() {
			continue
		}

		if !strings.Contains(fileName.Name(), ".") { // If the file name doesn't have a period in it - ignore!
			continue
		}

		var fileLine []string          // Get the lines of the file
		var filePath = fileName.Name() // Set the file name

		if slices.Contains(unwantedFiles, filePath) { // Make sure it's not one of the known unwanted files to edit
			continue
		}

		var unwantedExtention bool = false
		var updatedFile bool = false

		for _, extension := range unwantedExtentions { // ignore binary files!
			if strings.Contains(filePath, extension) {
				unwantedExtention = true
			}
		}

		if unwantedExtention {
			continue
		}

		file, err := os.Open(filePath) // Look for to dos in the file
		if err != nil {
			return
		}

		var lineNumber int
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			if strings.Contains(line, "TODO: ") && !strings.Contains(line, ") TODO") { // This is adding a number to the start of the todo as a way to keep track and act as a guard against duplicating issues!
				var replaceString string = fmt.Sprintf("(#%d) TODO", CurrentNumberOfIssues+1)
				line = strings.Replace(line, "TODO", replaceString, 1)
				fmt.Printf("I would like to make a github issue for: %s\nThe title is %s\nThe body is: %s on line %d\n", strings.TrimSpace(line), strings.TrimSpace(line), fileName.Name(), lineNumber)
				CurrentNumberOfIssues += 1
				// Check whether the issue already exists...
				MakeGithubIssue(line, fmt.Sprintf("This is from file %s on line %d\n", fileName.Name(), lineNumber))
				updatedFile = true
				foundNewTODO = true
			} else if strings.Contains(line, "TODO: ") && strings.Contains(line, ") TODO") {
				// This finds OLD TODOs

				// (#22) TODO: If github issue not in the list of old todos close issue
				/* _, removeError := RemoveLineDueToGithubIssue(line, listOfGithubIssues)
				if removeError == nil {
					// (#21) TODO: If todo in the list of old todos and no longer open on github, remove line
					line = ""
				} */

				// issue here being TOO powerful, when run on itself it deletes the if statements! Check for number?
			}
			fileLine = append(fileLine, line)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}

		// Write modified content back to the file
		if updatedFile {
			err = os.WriteFile(filePath, []byte(strings.Join(fileLine, "\n")), 0644)
			if err != nil {
				fmt.Println("Error writing file:", err)
				return
			}
		}
	}
	// MakeGithubIssue("My first GitHub Issue", "This is my first github issue from the API")

	if !foundNewTODO {
		fmt.Println("No new todo found in any file in this directory")
	}
}
