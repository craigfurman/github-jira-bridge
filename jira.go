package githubjirabridge

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/andygrunwald/go-jira"
	"github.com/google/go-github/v42/github"
)

var (
	issueStubRegexp = regexp.MustCompile(`(?m)^https://github.com/(?P<owner>[^/]+)/(?P<repo>[^/]+)/issues/(?P<number>\d+)\n\n` + jiraStubSignature + "$")
)

func handleJiraIssueWebhook(w http.ResponseWriter, req *http.Request) {
	logger.Println("handling jira issue webhook")

	token := req.URL.Query().Get("token")
	if token == "" {
		logger.Println("missing jira webhook token")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if subtle.ConstantTimeCompare([]byte(jiraWebhookToken), []byte(token)) == 0 {
		logger.Println("invalid jira webhook token")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var issueEvent JiraIssueEvent
	if err := json.NewDecoder(req.Body).Decode(&issueEvent); err != nil {
		logger.Printf("error parsing jira webhook: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Printf("received jira issue webhook for issue %s", issueEvent.Issue.ID)

	for _, changelogItem := range issueEvent.Changelog.Items {
		desiredStatusGithubLabel := jiraStatusToGithubLabels[changelogItem.ToString]
		if changelogItem.Field == "status" && desiredStatusGithubLabel != "" {
			logger.Printf("jira issue %s has changed status to %s, will ensure that related Github issue has label %s\n", issueEvent.Issue.Key, changelogItem.ToString, desiredStatusGithubLabel)
			handleStatusChangedJiraIssue(req.Context(), w, issueEvent, desiredStatusGithubLabel)
			return
		}
	}
	logger.Printf("will do nothing with jira issue %s\n", issueEvent.Issue.Key)
}

func handleStatusChangedJiraIssue(ctx context.Context, w http.ResponseWriter, issueEvent JiraIssueEvent, desiredStatusGithubLabel string) {
	owner, repo, number := parseGithubIssueFromJiraStub(issueEvent.Issue.Fields.Description)
	if number == 0 {
		logger.Printf("error parsing Github issue details from jira issue %s\n", issueEvent.Issue.Key)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ghIssue, _, err := githubClient.Issues.Get(ctx, owner, repo, number)
	if err != nil {
		logger.Printf("error getting issue #%d: %s\n", number, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	labels := ghIssue.Labels
	var newLabels []string
	for _, label := range labels {
		if !isGithubStatusLabel(*label.Name) && *label.Name != githubTriageLabel {
			newLabels = append(newLabels, *label.Name)
		}
	}
	newLabels = append(newLabels, desiredStatusGithubLabel)
	issuePatch := &github.IssueRequest{Labels: &newLabels}
	_, _, err = githubClient.Issues.Edit(ctx, owner, repo, number, issuePatch)
	if err != nil {
		logger.Printf("error updating issue #%d: %s\n", number, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type JiraIssueEvent struct {
	Issue     jira.Issue            `json:"issue"`
	Changelog jira.ChangelogHistory `json:"changelog"`
}

func parseGithubIssueFromJiraStub(description string) (string, string, int) {
	submatches := issueStubRegexp.FindStringSubmatch(description)
	var owner, repo, numberStr string
	for i, name := range issueStubRegexp.SubexpNames() {
		if name == "owner" && i < len(submatches) {
			owner = submatches[i]
		} else if name == "repo" && i < len(submatches) {
			repo = submatches[i]
		} else if name == "number" && i < len(submatches) {
			numberStr = submatches[i]
		}
	}
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return "", "", 0
	}
	return owner, repo, number
}

func isGithubStatusLabel(label string) bool {
	for _, githubLabel := range jiraStatusToGithubLabels {
		if label == githubLabel {
			return true
		}
	}
	return false
}
