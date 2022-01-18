package main

import (
	"net/http"
	"os"

	githubjirabridge "github.com/craigfurman/github-jira-bridge"
)

func main() {
	addr := os.Getenv("LISTEN_ADDRESS")
	if addr == "" {
		addr = "localhost:8080"
	}
	if err := http.ListenAndServe(addr, http.HandlerFunc(githubjirabridge.GitHubJiraBridge)); err != nil {
		panic(err)
	}
}
