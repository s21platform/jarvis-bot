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
