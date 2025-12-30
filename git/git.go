package git

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	utils "github.com/jonathon-chew/Thoth/Utils"
)

// (#24) TODO: Pretty printing for better reading?

// GITHUB STRUCTS
type Assignee struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type Issue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Milestone int      `json:"milestone,omitempty"`
	Label     []string `json:"labels,omitempty"`
	Assignees string   `json:"assignees,omitempty"`
}

type Label struct {
}

type GithubIssueResponse struct {
	Url            string `json:"url"`
	Repository_url string `json:"repository_url"`
	Labels_url     string `json:"labels_url"`
	Comments_url   string `json:"comments_url"`
	Events_url     string `json:"events_url"`
	Id             int    `json:"id"`
	Node_id        string `json:"node_id"`
	Number         int    `json:"number"`
	Title          string `json:"title"`
	User           struct {
		Login          string `json:"login"`
		Id             int    `json:"id"`
		Repos_url      string `json:"repos_url"`
		Events_url     string `json:"events_url"`
		Type           string `json:"type"`
		User_view_type string `json:"user_view_type"`
		Site_admin     bool   `json:"site_admin"`
	} `json:"user"`
	Labels             []Label    `json:"labels"`
	State              string     `json:"state"`
	State_Reason       string     `json:"state_reason"`
	Locked             bool       `json:"locked"`
	Assignee           Assignee   `json:"assignee"`
	Assignees          []Assignee `json:"assignees"`
	Comments           int        `json:"comments"`
	Created_at         string     `json:"created_at"`
	Updated_at         string     `json:"updated_at"`
	Author_association string     `json:"author_association"`
	Active_lock_reason string     `json:"active_lock_reason"`
	Body               string     `json:"body"`
	Message            string     `json:"message"`
	Status             string     `json:"status"`
}

var GithubStatusResponseMeanings = map[string]string{
	"201": "Created",
	"400": "Bad Request",
	"403": "Forbidden",
	"404": "Resource not found",
	"410": "Gone",
	"422": "Validation failed, or the endpoint has been spammed.",
	"503": "Service unavailable",
}

type Credentials struct {
	Owner string
	Repo  string
	Token string
}

// UTILS
func GetRemoteOrigin() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", stderr.String())
		return "", err
	}

	return out.String(), nil
}

func FindGitFolder() bool {

	directoryList := utils.MakeDirectoryList(utils.FindFilesInCurrentDirectory())

	// Look in the directories for a git folder
	if !slices.Contains(directoryList, ".git") {
		fmt.Println("[ERROR]: No git folder found")
		return false // recursively look?
	}

	return true
}

func OpenRemoteOrigin(place string) error {
	url, ErrGetRemote := GetRemoteOrigin()
	if ErrGetRemote != nil {
		return ErrGetRemote
	}

	url = strings.TrimSpace(url)

	if strings.Contains(url, "github.com") && place != "" {
		switch place {
		case "pull":
			url = url + "/pulls"
		case "issues":
			url = url + "/issues"
		}
	} else if place != "" {
		return fmt.Errorf("[ERROR]: only github.com has been implimented so far")
	}

	cmd := exec.Command("open", url)

	ErrRun := cmd.Run()
	if ErrRun != nil {
		fmt.Printf("Error: %s\n", ErrRun)
		return ErrRun
	}

	return nil
}

// GIT TAG
func GetTags() (string, error) {
	cmd := exec.Command("git", "tag")

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", stderr.String())
		return "", err
	}

	versions := out.String()

	return versions, nil
}

func GetLatestTag() (string, error) {

	if !FindGitFolder() {
		return "", fmt.Errorf("[Error]: Unable to find a git folder in the current directory")
	}

	versions, err := GetTags()
	if err != nil {
		return "", fmt.Errorf("[Error]: Unable to successfully get the tags\n ")
	}

	versionList := strings.Split(versions, "\n")

	// If the list is only 1 item long it's the biggest, so early return
	if len(versionList) == 1 {
		return versions, nil
	}

	var biggestMajor, biggestMinor, biggestPatch int
	possibleVersions := versionList

	for index, version := range versionList {

		if len(version) < 4 {
			continue
		}

		if !strings.Contains(version, ".") && !strings.Contains(version, "v") {
			fmt.Printf("[WARNING]: Skipping looking at tag %s, as doesn't follow the convention v.[0-9].[0-9].[0-9]", version)
			continue
			// return "", fmt.Errorf("[ERROR] unable to find periods in the version tags")
		}

		major, ErrMajorConv := strconv.Atoi(strings.Split(version[1:], ".")[0])
		minor, ErrMinorConv := strconv.Atoi(strings.Split(version[1:], ".")[1])
		patch, ErrPatchConv := strconv.Atoi(strings.Split(version[1:], ".")[2])

		if ErrMajorConv != nil && ErrMinorConv != nil && ErrPatchConv != nil {
			return "", fmt.Errorf("[ERROR]: There was an error converting %s, %s, %s", strings.Split(version[1:], ".")[0], strings.Split(version[1:], ".")[1], strings.Split(version[1:], ".")[2])
		}

		if major > biggestMajor {
			possibleVersions = versionList[index:]
		} else if major < biggestMajor {
			possibleVersions = append(possibleVersions[:index], possibleVersions[index+1:]...)
		}

		if len(possibleVersions) == 1 {
			return strings.Join(versionList, " "), nil
		}

		if minor > biggestMinor && slices.Contains(possibleVersions, version) {
			possibleVersions = versionList[index:]
		} else if minor < biggestMinor {
			possibleVersions = append(possibleVersions[:index], possibleVersions[index+1:]...)
		}

		if len(possibleVersions) == 1 {
			return strings.Join(versionList, " "), nil
		}

		if patch > biggestPatch {
			possibleVersions = versionList[index:]
		} else if patch < biggestPatch {
			possibleVersions = append(possibleVersions[:index], possibleVersions[index+1:]...)
		}

		if len(possibleVersions) == 1 {
			return strings.Join(versionList, " "), nil
		}
	}

	return strings.Join(possibleVersions, ""), nil
}

func NewGitTag(argument string) error {
	version, ErrGetLatestTag := GetLatestTag()
	if ErrGetLatestTag != nil {
		return ErrGetLatestTag
	}
	fmt.Println("Current latest tag: ", version)

	if argument != "major" && argument != "minor" && argument != "patch" {
		var userChoiceVersionUpdate string

		fmt.Printf("Do you want to increase the major, minor or patch of the tag?\n")

		_, ErrUserInput := fmt.Scanln(&userChoiceVersionUpdate)
		if ErrUserInput != nil {
			return ErrUserInput
		}
		if userChoiceVersionUpdate != "major" && userChoiceVersionUpdate != "minor" && userChoiceVersionUpdate != "patch" {
			return fmt.Errorf("[ERROR]: user input was not major, minor or patch")
		} else {
			argument = userChoiceVersionUpdate
		}
	}

	major, ErrMajorConv := strconv.Atoi(strings.Split(version[1:], ".")[0])
	if ErrMajorConv != nil {
		return ErrMajorConv
	}

	minor, ErrMinorConv := strconv.Atoi(strings.Split(version[1:], ".")[1])
	if ErrMinorConv != nil {
		return ErrMinorConv
	}

	patch, ErrPatchConv := strconv.Atoi(strings.Split(version[1:], ".")[2])
	if ErrPatchConv != nil {
		return ErrPatchConv
	}
	var newTag string

	switch argument {
	case "major":
		newMajor := major + 1
		newTag = fmt.Sprintf("v%d.%d.%d", newMajor, 0, 0)
	case "minor":
		newMinor := minor + 1
		newTag = fmt.Sprintf("v%d.%d.%d", major, newMinor, 0)
	case "patch":
		newPatch := patch + 1
		newTag = fmt.Sprintf("v%d.%d.%d", major, minor, newPatch)
	}

	cmd := exec.Command("git", "tag", newTag)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", stderr.String())
		return err
	}

	aphrodite.PrintInfo(fmt.Sprintf("New latest tag:%s\n", newTag))

	aphrodite.PrintBold("Cyan", "Do you want to push the new tag to git?\n")

	var userChoicePushToGit string
	_, ErrGettingUserChioce := fmt.Scan(&userChoicePushToGit)
	if ErrGettingUserChioce != nil {
		return ErrGettingUserChioce
	}

	if userChoicePushToGit == "y" || userChoicePushToGit == "Y" || userChoicePushToGit == "yes" || userChoicePushToGit == "Yes" || userChoicePushToGit == "YES" {
		aphrodite.PrintInfo("Pushing to remote git respository.\n")
		// git push --tags --force-with-lease=false
		tagPushCmd := exec.Command("git", "push", "--tags", "--force-with-lease=false")
		ErrPushingTags := tagPushCmd.Run()
		if ErrPushingTags != nil {
			return ErrPushingTags
		}
		aphrodite.PrintInfo("Successfully pushed.\n")
	}

	return nil
}

// MAKE A GIT REQUEST
func GenericGitRequest() (Credentials, error) {
	remoteOrigin, err := GetRemoteOrigin()
	var credentials Credentials
	if err != nil {
		fmt.Printf("Unable to get the remote origin\n")
		return credentials, err
	}

	if strings.Contains(remoteOrigin, "github") {
		gitUrl := strings.ReplaceAll(remoteOrigin, ".git", "")
		gitDetails := strings.Split(strings.ReplaceAll(gitUrl, "https://github.com/", ""), "/")

		credentials.Owner = gitDetails[0]
		credentials.Repo = strings.Replace(gitDetails[1], "\n", "", -1)
		credentials.Token = os.Getenv("GH_PERSONAL_TOKEN")

		if credentials.Token == "" {
			return credentials, errors.New("no GH_PERSONAL_TOKEN in the environment")
		}

		return credentials, nil
	} else {
		return credentials, fmt.Errorf("the remote origin is not github, and the ability to create issues for %s is not currently implimented", remoteOrigin)
	}
}

// LIST GIT ISSUES
func ListGithubIssues(passedFromCLI bool) ([]GithubIssueResponse, error) {

	var ResponseInstance []GithubIssueResponse

	GitCredentials, err := GenericGitRequest()
	if err != nil {
		return ResponseInstance, err
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=all", GitCredentials.Owner, GitCredentials.Repo), nil)
	if err != nil {
		return ResponseInstance, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", GitCredentials.Token))

	client := http.Client{}

	req, err := client.Do(request)
	if err != nil {
		return ResponseInstance, err
	}

	defer req.Body.Close()

	// (#26) TODO: Decide on whether or not you want this printed out EVERY time the programme is run

	if !passedFromCLI {
		fmt.Printf("The response was: %s, %s\n\n", req.Status, GithubStatusResponseMeanings[req.Status])
	}

	responseBody, err := io.ReadAll(req.Body)
	if err != nil {
		return ResponseInstance, err
	}

	// fmt.Printf("Repsonse Body: %s\n\n", string(responseBody))

	if err := json.Unmarshal(responseBody, &ResponseInstance); err != nil {
		return ResponseInstance, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if len(ResponseInstance) == 0 {
		return ResponseInstance, errors.New("no GitHub issues found")
	}

	if req.StatusCode != http.StatusOK {
		return ResponseInstance, fmt.Errorf("GitHub API error: %s", req.Status)
	}

	// fmt.Printf("ResponseInstance: %v\n\n", ResponseInstance)

	// for _, response := range ResponseInstance {
	// 	fmt.Println("The title for the response is: ", strings.TrimSpace(response.Title), " with ID: ", response.Id)
	// }

	return ResponseInstance, nil
}

func MakeGithubIssue(TITLE, BODY string) error {

	// Get the credentials required
	GithubCredentials, err := GenericGitRequest()
	if err != nil {
		return err
	}

	// Create the issue using a struct
	issue := Issue{
		Title: strings.TrimSpace(TITLE),
		Body:  BODY,
	}

	// Convert the struct into JSON using the tags and Marshal
	jsonData, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	// Convert the JSON into bytes
	requestBody := bytes.NewBuffer(jsonData)

	// Make the request
	request, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", GithubCredentials.Owner, GithubCredentials.Repo), io.Reader(requestBody))
	if err != nil {
		fmt.Printf("Error making the HTTP request %s\n", err)
		return err
	}

	// Set the required headers
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", GithubCredentials.Token))

	// Make a new client
	client := http.Client{}

	// Complete the request - Client.Do because the http.NewRequest handles the method
	req, err := client.Do(request)
	if err != nil {
		return err
	}

	if req.StatusCode != 200 && req.StatusCode != 201 {
		fmt.Println(req.Body)
		return fmt.Errorf("the response was not positive, %d", req.StatusCode)
	}

	fmt.Printf("The response was: %s, %s\n", req.Status, GithubStatusResponseMeanings[req.Status])

	return nil
}

// REMOVE GIT ISSUES
// (#2) TODO: Add the ability to remove to dos which have been closed on github
func RemoveLineDueToGithubIssue(line string, foundGithubIssues []GithubIssueResponse) (bool, error) {

	// Loop through the issues and compare to the line
	for _, issue := range foundGithubIssues {
		if strings.Contains(strings.TrimSpace(line), issue.Title) {
			err := CloseGithubIssue(&issue)
			if err != nil {
				return true, err // trying this out - as first half the of the function was "completed" successfully but the second half wasn't!
			}
			return true, nil
		}
	}

	// If the loop didn't find anything return false and no error!
	return false, nil
}

// (#3) TODO: Add the ability to close issues on github which have been removed from the code base

func CloseGithubIssue(closeIssue *GithubIssueResponse) error {

	// Put together the JSON message required to close an issue
	closeIssue.State = "closed"
	closeIssue.State_Reason = "completed"

	// Get the credentials
	GithubCredentials, err := GenericGitRequest()
	if err != nil {
		return err
	}

	// Convert the struct into JSON using the tags and Marshal
	jsonData, err := json.Marshal(closeIssue)
	if err != nil {
		return err
	}

	// Convert the JSON into bytes
	requestBody := bytes.NewBuffer(jsonData)

	// Write the request
	request, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", GithubCredentials.Owner, GithubCredentials.Repo, closeIssue.Number), requestBody)
	if err != nil {
		return err
	}

	// Set the required headers
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", GithubCredentials.Token))

	client := http.Client{}

	// Make the request
	closeGithubIssueResponse, clientErr := client.Do(request)
	if clientErr != nil {
		fmt.Printf("The response from github was: %s\n", GithubStatusResponseMeanings[closeGithubIssueResponse.Status])
		return clientErr
	}

	fmt.Printf("The response from github was: %s\n", GithubStatusResponseMeanings[closeGithubIssueResponse.Status])

	// Return if error?
	return nil

}
