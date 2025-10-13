package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	"github.com/jonathon-chew/Thoth/git"
)

func CLI(CommandLineArguments []string) error {
	aphrodite.PrintColour("Cyan", "I have found additional command line arguments, switching to CLI mode\n")

	var NoIssues error = errors.New("no GitHub issues found")

	for index, command := range CommandLineArguments {
		switch command {
		case "--get", "-get", "-g":
			returned, err := git.ListGithubIssues(true)
			if err != nil && errors.Is(err, NoIssues) {
				aphrodite.PrintWarning("no GitHub issues found")
				return nil
			}

			if err != nil {
				return err
			}

			var closedFlag, openFlag bool = false, false
			// Check for extra flags
			if len(os.Args) > 2 {
				for _, extraCommand := range os.Args[2:] {
					switch extraCommand {
					case "--closed", "-closed", "-c":
						closedFlag = true
					case "--open", "-open", "-o":
						openFlag = true
					}
				}
			}

			for index, issue := range returned {
				if closedFlag && issue.State == "closed" {
					fmt.Printf("%d The issue title is:%s\nThe body is:%s\nThe status is:%s\n\n", index+1, strings.TrimSpace(issue.Title), issue.Body, aphrodite.ReturnInfo(issue.State))
					continue
				}

				if openFlag && issue.State == "open" {
					fmt.Printf("%d The issue title is:%s\nThe body is:%s\nThe status is:%s\n\n", index+1, strings.TrimSpace(issue.Title), issue.Body, aphrodite.ReturnError(issue.State))
					continue
				}

				if !closedFlag && !openFlag {
					fmt.Printf("%d The issue title is:%s\nThe body is:%s\nThe status is:%s\n\n", index+1, strings.TrimSpace(issue.Title), issue.Body, issue.State)
				}
			}

			return nil

		case "--set", "-set", "-s":
			var IssueTitle, IssueBody string
			if CommandLineArguments[index+1] == "title" || CommandLineArguments[index+1] == "--title" || CommandLineArguments[index+1] == "-title" || CommandLineArguments[index+1] == "-t" {
				IssueTitle = CommandLineArguments[index+2]
			} else {
				return errors.New("could not find a title flag proceeding the set command")
			}

			if CommandLineArguments[index+3] == "body" || CommandLineArguments[index+3] == "--body" || CommandLineArguments[index+3] == "-body" || CommandLineArguments[index+3] == "-b" {
				IssueBody = CommandLineArguments[index+4]
			} else {
				return errors.New("could not find a body flag proceeding the set command")
			}

			makeError := git.MakeGithubIssue(IssueTitle, IssueBody)
			if makeError != nil {
				fmt.Println(makeError)
				return makeError
			}

			return nil
		case "--version", "-version", "-v":
			fmt.Printf("v0.0.4\n")
		case "--help", "-help", "-h":
			type Help struct {
				NoArguments string
				GetIssues   string
				SetIssues   string
				Version     string
			}

			var newHelp Help

			newHelp.NoArguments, _ = aphrodite.ReturnBold("Cyan", "No Arguments")
			newHelp.GetIssues, _ = aphrodite.ReturnBold("Cyan", "Get issues")
			newHelp.SetIssues, _ = aphrodite.ReturnBold("Cyan", "Set issues")
			newHelp.Version, _ = aphrodite.ReturnBold("Cyan", "Version")

			fmt.Printf("\n%s\nYou can run with no arguments to check all files\n%s\nYou can pass in a get flag which will List the github issues\n%s\nIf you pass in the set flag, please pass in the title flag and body flag (in that order) to make a new issue with the relevent Title and Body\n%s\nVersion Number can be passed in with the version flag", newHelp.NoArguments, newHelp.GetIssues, newHelp.SetIssues, newHelp.Version)
		}

	}
	return nil
}
