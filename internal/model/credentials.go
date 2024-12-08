package model

import "encoding/json"

type CredData struct {
	Data json.RawMessage `db:"data"`
}

type Data struct {
	Creds []Cred `json:"creds"`
}

type Cred struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
