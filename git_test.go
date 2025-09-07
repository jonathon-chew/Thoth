package main

import (
	"testing"
)

func TestRemoteURL(t *testing.T) {
	t.Logf("Testing GetRemoteOrigin")
	url, err := GetRemoteOrigin()
	if err != nil {
		t.Fatalf("Failed to get remote origin: %v", err)
	}

	t.Logf("Remote URL: %s", url)
}

func TestListIssues(t *testing.T) {
	t.Logf("Testing ListGithubIssues")
	returned, err := ListGithubIssues()
	if err != nil {
		t.Fatalf("Failed to get remote origin: %v", err)
	}

	for _, i := range returned{
		t.Logf("From the test file: The title from the API was: %v\n", i.Title)
	}

}

func TestGenericGit(t *testing.T) {
	t.Logf("Testing GetRemoteOrigin")
	GitCredentials, err := genericGitRequest()
	if err != nil {
		t.Fatalf("Failed to get Git data: %v", err)
	}

	t.Logf("Owner: %s, Repo: %s, Token: %s", GitCredentials.Owner, GitCredentials.Repo, GitCredentials.Token)
}