package models

import "time"

type Todo struct {
	Id          string     `json:"id,omitempty"`
	UserId      string     `json:"user_id,omitempty"`
	Deadline    *time.Time `json:"deadline"`
	Description string     `json:"description"`
}

func (t Todo) GetMapping() string {
	return `
	{
		"mappings": {
		"_doc": {
			"properties": {
				"user_id": {
					"type": "keyword"
				},
				"deadline": {
					"type": "date"
				},
				"description": {
					"type": "text"
				}
			}
		}
	}
	}
`
}
