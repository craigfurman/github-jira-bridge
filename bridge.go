package githubjirabridge

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/google/go-github/v42/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

const (
	jiraStubSignature = "Created by github-jira-bridge."
)

var (
	router                   *mux.Router
	logger                   *log.Logger
	jiraClient               *jira.Client
	githubClient             *github.Client
	jiraStatusToGithubLabels map[string]string
)

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("starting up")

	router = mux.NewRouter()
	router.HandleFunc("/github", handleGithubWebhook).Methods(http.MethodPost)
	router.HandleFunc("/jira/issues", handleJiraIssueWebhook).Methods(http.MethodPost)

	jiraAuth := jira.BasicAuthTransport{
		Username: jiraUsername,
		Password: jiraApiToken,
	}
	var err error
	jiraClient, err = jira.NewClient(jiraAuth.Client(), jiraApiURL)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient = github.NewClient(tc)

	jiraStatusToGithubLabels = map[string]string{}
	for _, pair := range strings.Split(jiraStatusChangeToGithubLabels, ",") {
		if pair == "" {
			continue
		}
		jiraStatusAndGithubLabel := strings.Split(pair, ":")
		jiraStatusToGithubLabels[jiraStatusAndGithubLabel[0]] = jiraStatusAndGithubLabel[1]
	}
}

// Google Cloud Function entrypoint
func GitHubJiraBridge(w http.ResponseWriter, req *http.Request) {
	router.ServeHTTP(w, req)
}
