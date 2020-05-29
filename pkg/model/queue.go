package model

type Queue struct {
	Queue   string   `json:"queue"`
	Members []string `json:"members"`
}
