package githubjirabridge

import "os"

var (
	jiraApiURL     = os.Getenv("JIRA_API_URL")
	jiraUsername   = os.Getenv("JIRA_USERNAME")
	jiraApiToken   = os.Getenv("JIRA_API_TOKEN")
	jiraProjectKey = os.Getenv("JIRA_PROJECT_KEY")
	jiraIssueType  = os.Getenv("JIRA_ISSUE_TYPE")

	githubWebhookSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")
	githubTriageLabel   = os.Getenv("GITHUB_TRIAGE_LABEL")
)
