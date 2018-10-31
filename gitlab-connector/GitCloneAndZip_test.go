package main

import (
	"testing"
)

func TestGitCloneAndZip(t *testing.T) {
	t.Run("Repository not available", func(t *testing.T) {
		testData := &GitCloneAndZipInput{repositoryUrl: "https://0.0.0.0:12345/gitlab/testrepo.git"}

		_, err := GitCloneAndZip(testData)
		if err == nil {
			t.Fatal("Error: Should fail on unavailable repositories")
		}
	})

}
