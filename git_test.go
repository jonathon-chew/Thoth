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
	returned, err := ListGithubIssues(false) // false is NOT passed from the CLI so will always report if it connected to github
	if err != nil {
		t.Fatalf("Failed to get remote origin: %v", err)
	} else {
		t.Logf("Found issues %d", len(returned))
	}

	// for _, i := range returned{
	// 	t.Logf("From the test file: The title from the API was: %v\n", i.Title)
	// }

}

func TestGenericGit(t *testing.T) {
	t.Logf("Testing GetRemoteOrigin")
	GitCredentials, err := genericGitRequest()
	if err != nil {
		t.Fatalf("Failed to get Git data: %v", err)
	}

	t.Logf("Owner: %s, Repo: %s, Token: %s", GitCredentials.Owner, GitCredentials.Repo, GitCredentials.Token)
}
