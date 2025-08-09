package jira

import (
	"log"

	"github.com/andygrunwald/go-jira"

	"github.com/s21platform/jarvis-bot/internal/config"
)

type Jira struct {
	*jira.Client
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
		client,
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

func (j *Jira) CreateIssue(title string, issueType string, projectKey string) (string, error) {
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary: title,
			Type: jira.IssueType{
				Name: issueType,
			},
			Project: jira.Project{
				Key: projectKey,
			},
		},
	}

	resp, _, err := j.Issue.Create(issue)
	if err != nil {
		return "", err
	}
	return resp.Key, nil
}
