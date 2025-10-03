package main

import (
	"errors"
	"fmt"
)

func CLI(CommandLineArguments []string) error {
	fmt.Printf("I have found additional command line arguments, switching to CLI mode\n\n")

	for index, command := range CommandLineArguments {
		switch command {
		case "--get", "-get", "-g":
			returned, err := ListGithubIssues()
			if err != nil {
				return err
			}

			for _, issue := range returned {
				fmt.Printf("The issue title is:%v\n", issue.Title)
			}

			return nil

		case "--set", "-set", "-s":
			var IssueTitle, IssueBody string
			if CommandLineArguments[index+1] == "title" || CommandLineArguments[index+1] == "--title" || CommandLineArguments[index+1] == "-title" {
				IssueTitle = CommandLineArguments[index+2]
			} else {
				return errors.New("could not find a title flag proceeding the set command")
			}

			if CommandLineArguments[index+3] == "body" || CommandLineArguments[index+3] == "--body" || CommandLineArguments[index+3] == "-body" {
				IssueBody = CommandLineArguments[index+4]
			} else {
				return errors.New("could not find a body flag proceeding the set command")
			}

			makeError := MakeGithubIssue(IssueTitle, IssueBody)
			if makeError != nil {
				fmt.Println(makeError)
				return makeError
			}

			return nil
		case "--version", "-version", "-v":
			fmt.Printf("v0.0.2\n")
		case "--help", "-help", "-h":
			fmt.Printf("\nNo Arguments\nYou can run with no arguments to check all files\nGet Issues\nYou can pass in a get flag which will List the github issues\nNew Issue\nIf you pass in the set flag, please pass in the title flag and body flag (in that order) to make a new issue with the relevent Title and Body\nVersion Number\nVersion Number can be passed in with the version flag")
		}

	}
	return nil
}
