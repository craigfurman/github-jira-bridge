package githubjirabridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGithubIssueFromJiraStub(t *testing.T) {
	for _, tc := range []struct {
		description         string
		expectedOwner       string
		expectedRepo        string
		expectedIssueNumber int
	}{
		{
			description:         "https://github.com/an-org/a-repo/issues/1\n\nCreated by github-jira-bridge.",
			expectedOwner:       "an-org",
			expectedRepo:        "a-repo",
			expectedIssueNumber: 1,
		},
		{
			description:         "https://github.com/an-org/a-repo/issues/10\n\nCreated by github-jira-bridge.",
			expectedOwner:       "an-org",
			expectedRepo:        "a-repo",
			expectedIssueNumber: 10,
		},
		{
			description:         "A random issue.",
			expectedIssueNumber: 0,
		},
		{
			description:         "https://github.com/an-org/a-repo/issues/123\n\nThis isn't really a stub.",
			expectedIssueNumber: 0,
		},
		{
			description:         "Some text\n\nhttps://github.com/an-org/a-repo/issues/123\n\nThis isn't really a stub.",
			expectedIssueNumber: 0,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			owner, repo, issueNumber := parseGithubIssueFromJiraStub(tc.description)
			assert.Equal(t, tc.expectedOwner, owner)
			assert.Equal(t, tc.expectedRepo, repo)
			assert.Equal(t, tc.expectedIssueNumber, issueNumber)
		})
	}
}
