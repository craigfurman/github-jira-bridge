package githubjirabridge

import (
	"log"
	"net/http"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/gorilla/mux"
)

var (
	router     *mux.Router
	logger     *log.Logger
	jiraClient *jira.Client
)

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("starting up")

	router = mux.NewRouter()
	router.HandleFunc("/github", handleGithubWebhook).Methods(http.MethodPost)

	jiraAuth := jira.BasicAuthTransport{
		Username: jiraUsername,
		Password: jiraApiToken,
	}
	var err error
	jiraClient, err = jira.NewClient(jiraAuth.Client(), jiraApiURL)
	if err != nil {
		panic(err)
	}
}

// Google Cloud Function entrypoint
func GitHubJiraBridge(w http.ResponseWriter, req *http.Request) {
	router.ServeHTTP(w, req)
}
