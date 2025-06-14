package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

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
	Labeles            []string   `json:"labels"`
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

var GithubStatusResponseMeanings = map[string]string{
	"201": "Created",
	"400": "Bad Request",
	"403": "Forbidden",
	"404": "Resource not found",
	"410": "Gone",
	"422": "Validation failed, or the endpoint has been spammed.",
	"503": "Service unavailable",
}

func genericGitRequest() (OWNER string, REPO string, TOKEN string, err error) {
	remoteOrigin, err := GetRemoteOrigin()
	if err != nil {
		fmt.Printf("Unable to get the remote origin\n")
		return "", "", "", err
	}

	if strings.Contains(remoteOrigin, "github") {
		// https://github.com/OWNER/REPO.git
		gitUrl := strings.ReplaceAll(remoteOrigin, ".git", "")
		gitDetails := strings.Split(strings.ReplaceAll(gitUrl, "https://github.com/", ""), "/")

		OWNER = gitDetails[0]
		REPO = strings.Replace(gitDetails[1], "\n", "", -1)
		TOKEN = os.Getenv("GH_PERSONAL_TOKEN")

		return OWNER, REPO, TOKEN, nil
	} else {
		return "", "", "", errors.New(fmt.Sprintf("The remote origin is not github, and the ability to create issues for %s is not currently implimented.", remoteOrigin))
	}
}

func ListGithubIssues() ([]GithubIssueResponse, error) {

	var ResponseInstance []GithubIssueResponse

	OWNER, REPO, TOKEN, err := genericGitRequest()
	if err != nil {
		return ResponseInstance, err
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", OWNER, REPO), nil)
	if err != nil {
		return ResponseInstance, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", TOKEN))

	client := http.Client{}

	req, err := client.Do(request)
	if err != nil {
		return ResponseInstance, err
	}

	defer req.Body.Close()

	// fmt.Printf("The response was: %s, %s\n\n", req.Status, GithubStatusResponseMeanings[req.Status])

	responseBody, err := io.ReadAll(req.Body)
	if err != nil {
		return ResponseInstance, err
	}

	// fmt.Printf("Repsonse Body: %s\n\n", string(responseBody))

	if string(responseBody) == "[]" {
		CustomResponseError := fmt.Errorf("There were no github issues")
		return ResponseInstance, CustomResponseError
	}

	issues := json.Unmarshal(responseBody, &ResponseInstance)
	if issues != nil {
		fmt.Printf("Error Unmarshalling, %v\n", issues)
		return ResponseInstance, issues
	}

	// (#3) TODO: Check to see if this works properly, testing returns wrong?!
	if ResponseInstance[0].Status != "200" && len(ResponseInstance[0].Status) > 0 {
		CustomResponseError := fmt.Errorf("There was an error getting the github issues, %s\n", ResponseInstance[0].Message)
		return ResponseInstance, CustomResponseError
	}

	// fmt.Printf("ResponseInstance: %v\n\n", ResponseInstance)

	for _, i := range ResponseInstance {
		fmt.Println("The title for the response is: ", strings.TrimSpace(i.Title), " with ID: ", i.Id)
	}

	return ResponseInstance, nil
}

func MakeGithubIssue(TITLE, BODY string) error {

	// Get the credentials required
	OWNER, REPO, TOKEN, err := genericGitRequest()
	if err != nil {
		return err
	}

	// Create the issue using a struct
	issue := Issue{
		Title: TITLE,
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
	request, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", OWNER, REPO), io.Reader(requestBody))
	if err != nil {
		fmt.Printf("Error making the HTTP request %s\n", err)
		return err
	}

	// Set the required headers
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", TOKEN))

	// Make a new client
	client := http.Client{}

	// Complete the request - Client.Do because the http.NewRequest handles the method
	req, err := client.Do(request)
	if err != nil {
		return err
	}

	fmt.Printf("The response was: %s, %s\n", req.Status, GithubStatusResponseMeanings[req.Status])

	// defer req.Body.Close()

	return nil
}

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
	closeIssue.State        = "closed"
	closeIssue.State_Reason = "completed"

	// Get the credentials
	OWNER, REPO, TOKEN, err := genericGitRequest()
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
	request, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", OWNER, REPO, closeIssue.Number), requestBody)
	if err != nil {
		return err
	}

	// Set the required headers
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", TOKEN))

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