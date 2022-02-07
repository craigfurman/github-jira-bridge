package githubjirabridge

import "os"

var (
	jiraApiURL       = os.Getenv("JIRA_API_URL")
	jiraUsername     = os.Getenv("JIRA_USERNAME")
	jiraApiToken     = os.Getenv("JIRA_API_TOKEN")
	jiraProjectKey   = os.Getenv("JIRA_PROJECT_KEY")
	jiraIssueType    = os.Getenv("JIRA_ISSUE_TYPE")
	jiraWebhookToken = os.Getenv("JIRA_WEBHOOK_TOKEN")

	// Must have repo scope
	githubToken = os.Getenv("GITHUB_API_TOKEN")

	githubWebhookSecret   = os.Getenv("GITHUB_WEBHOOK_SECRET")
	githubTriageLabel     = os.Getenv("GITHUB_TRIAGE_LABEL")
	triagedIssueJiraLabel = os.Getenv("TRIAGED_ISSUE_JIRA_LABEL")

	// Format: jiraStatus1:githubLabel1,jiraStatus2:githubLabel2
	// Example: "To Do:status/todo,In Progress:status/in-progress"
	jiraStatusChangeToGithubLabels = os.Getenv("JIRA_STATUS_CHANGE_TO_GITHUB_LABEL_MAPPINGS")
)
