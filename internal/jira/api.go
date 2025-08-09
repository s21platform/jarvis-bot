package jira

import (
	"log"

	"github.com/andygrunwald/go-jira"

	"github.com/s21platform/jarvis-bot/internal/config"
)

type Jira struct {
	*jira.Client
	baseURL string
}

func New(cfg *config.Config) *Jira {
	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.Username,
		Password: cfg.Jira.Password,
	}

	client, err := jira.NewClient(tp.Client(), cfg.Jira.BaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	return &Jira{
		Client:  client,
		baseURL: cfg.Jira.BaseUrl,
	}
}

func (j *Jira) GetProjects() {
	projects, _, err := j.Project.GetList()
	if err != nil || projects == nil {
		log.Fatal(err)
	}
	for _, i := range *projects {
		log.Println(i.Key, i.Name)
	}
}

// CreateIssue создает новую задачу в Jira
func (j *Jira) CreateIssue(title string, issueType string, projectKey string) (string, error) {
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Project: jira.Project{
				Key: projectKey,
			},
			Summary: title,
			Type: jira.IssueType{
				Name: issueType,
			},
		},
	}

	issue, _, err := j.Issue.Create(&i)
	if err != nil {
		return "", err
	}

	return issue.Key, nil
}

// GetBaseURL возвращает базовый URL Jira
func (j *Jira) GetBaseURL() string {
	return j.baseURL
}
