package githubjirabridge

import (
	"fmt"
	"io"
	"net/http"

	"github.com/andygrunwald/go-jira"
	"github.com/google/go-github/v42/github"
)

const (
	jiraLabelForLinkedGitHubIssue = "from-github"
)

func handleGithubWebhook(w http.ResponseWriter, req *http.Request) {
	logger.Println("handling github webhook")
	payload, err := github.ValidatePayload(req, []byte(githubWebhookSecret))
	if err != nil {
		logger.Printf("error validating github payload: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(req), payload)
	if err != nil {
		logger.Printf("error parsing github webhook: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch event := event.(type) {
	case *github.IssuesEvent:
		handleGithubIssueEvent(w, event)
	default:
		logger.Printf("unrecognised GitHub event type: %s\n", req.Header.Get(github.EventTypeHeader))
	}
}

func handleGithubIssueEvent(w http.ResponseWriter, issue *github.IssuesEvent) {
	if *issue.Action == "labeled" && *issue.Label.Name == githubTriageLabel {
		logger.Printf("GitHub issue #%d was labeled for triage, will ensure there exists a matching jira issue\n", *issue.Issue.Number)
		ensureGithubIssueHasLinkedJiraIssue(w, issue)
		return
	}
	logger.Printf("GitHub issue #%d was %s, will do nothing\n", *issue.Issue.Number, *issue.Action)
}

// TODO make idempotent
func ensureGithubIssueHasLinkedJiraIssue(w http.ResponseWriter, issue *github.IssuesEvent) {
	logger.Printf("creating jira issue for GitHub issue #%d\n", *issue.Issue.Number)
	jiraIssueBody := fmt.Sprintf("%s\n\nCreated by github-jira-bridge.", *issue.Issue.HTMLURL)
	jiraIssue := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     *issue.Issue.Title,
			Description: jiraIssueBody,
			Labels:      []string{jiraLabelForLinkedGitHubIssue},
			Type:        jira.IssueType{Name: jiraIssueType},
			Project:     jira.Project{Key: jiraProjectKey},
		},
	}
	_, resp, err := jiraClient.Issue.Create(jiraIssue)
	if err != nil {
		var jiraRespBody []byte
		if resp != nil {
			defer resp.Body.Close()
			jiraRespBody, _ = io.ReadAll(resp.Body)
		}
		logger.Printf("error creating jira issue: %s: %s\n", err, string(jiraRespBody))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
