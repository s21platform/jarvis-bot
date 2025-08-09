package command

type JiraClient interface {
	CreateIssue(title string, issueType string, projectKey string) (string, error)
	GetBaseURL() string
}
