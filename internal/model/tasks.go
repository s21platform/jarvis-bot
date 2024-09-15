package model

type TasksByUUID struct {
	ID              int64  `json:"id"`
	TaskType        string `db:"task_type"`
	TaskTitle       string `db:"task_title"`
	TaskDescription string `db:"task_description"`
}
