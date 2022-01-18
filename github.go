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

func ensureGithubIssueHasLinkedJiraIssue(w http.ResponseWriter, issue *github.IssuesEvent) {
	signature := "Created by github-jira-bridge."
	jqlFindPotentiallyLinkedIssues := fmt.Sprintf(
		`project = %s AND labels = %s AND text ~ "%s\n" AND text ~ "%s"`,
		jiraProjectKey, jiraLabelForLinkedGitHubIssue, *issue.Issue.HTMLURL, signature,
	)
	potentiallyLinkedJiraIssues, resp, err := jiraClient.Issue.Search(jqlFindPotentiallyLinkedIssues, nil)
	if err != nil {
		var jiraRespBody []byte
		if resp != nil {
			defer resp.Body.Close()
			jiraRespBody, _ = io.ReadAll(resp.Body)
		}
		logger.Printf("error searching for jira issues: %s: %s\n", err, string(jiraRespBody))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(potentiallyLinkedJiraIssues) == 1 {
		logger.Printf("found jira issue %s that is linked to GitHub issue #%d, will do nothing\n", potentiallyLinkedJiraIssues[0].Key, *issue.Issue.Number)
		return
	}
	if len(potentiallyLinkedJiraIssues) > 1 {
		logger.Printf("found %d jira issues that are linked to GitHub issue #%d, will do nothing, but this is very unexpected\n", len(potentiallyLinkedJiraIssues), *issue.Issue.Number)
		return
	}

	logger.Printf("creating jira issue for GitHub issue #%d\n", *issue.Issue.Number)
	jiraIssueBody := fmt.Sprintf("%s\n\n%s", *issue.Issue.HTMLURL, signature)
	jiraIssue := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     *issue.Issue.Title,
			Description: jiraIssueBody,
			Labels:      []string{jiraLabelForLinkedGitHubIssue},
			Type:        jira.IssueType{Name: jiraIssueType},
			Project:     jira.Project{Key: jiraProjectKey},
		},
	}
	_, resp, err = jiraClient.Issue.Create(jiraIssue)
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
