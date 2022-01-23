# github-jira-bridge

Creates jira issues in reaction to certain GitHub issue event webhooks.

## Use cases

Intended users are teams that use jira internally to organise work, and also
maintain open-source projects on GitHub. This is a small project built to
scratch my own itch, but in the unlikely event other people discover this, I'm
open to contributions. I'm also open to new name suggestions for the project.

A few notes about the workflow this project aims to support:

1. Working in the open on open-source projects. Issues and discussions are
   public.
1. Using jira internally to manage issues pertaining to closed-source projects.
1. Using jira internally as a planning and prioritisation tool, across several
   projects.
1. Community management work is still work, and work needs to be prioritised.
   New public issues raised on open-source projects are welcome, but should not
   be automatically inserted into the backlog until triaged.

To support this, this project has the following features:

1. In response to a triage label being applied to a GitHub issue, create a
   linked stub jira issue.
1. In response to a jira issue's state changing, change a mutually exclusive
   status label on a linked GitHub issue.
1. **Coming soon**: In response to a GitHub issue being closed, change a linked
   jira issue's status to done.

At the moment, each instance of this project supports only 1 github and jira
project pair. I'm not against extending it to support multiple project mappings
per instance, but I personally don't need this yet.

## Usage

### Deployment

This project is designed to be deployed as a Google Cloud Function. That's why
the code is structured in a rather unusual way, with an `http.HandlerFunc` in
the root package. This is my first time using cloud functions and I only skimmed
the docs, so perhaps there's another way to structure it while still being able
to use them.

A main package is provided under `./cmd/server/`, which I currently use for
testing.

```
gcloud functions deploy GitHubJiraBridge --runtime go116 --trigger-http
```

Configure all of the environment variables found in [`config.go`](./config.go),
according to your particular jira/github projects.



Configure a github webhook for issues events, with the URL set to
`${CLOUD_FUNCTION_BASE_URL}/github`. Generate a random webhook secret and
configure this in both the webhook, and the relevent environment variable.

Configure a jira webhook for issues events, with the URL set to
`${CLOUD_FUNCTION_BASE_URL}/jira?token=abcd`. Set the JQL filter to `labels =
from-github`. Generate a random webhook secret and configure this in both the
token query parameter, and the relevent environment variable.

## Contributing

Contributions welcome, provided they're aligned with the general use-case. If in
doubt, please open an issue.

The code was hastily written in one evening, with no tests. If this project ends
up growing, that will likely need to change.

At the time of writing the logs are unstructured debug/trace style logs,
designed to be read linearly. I normally favour event-style logs, but given the
hasty and non-prod nature of this project I didn't want to spend the effort on
proper instrumentation middleware. If the project does grow, this is another
thing that might change.
