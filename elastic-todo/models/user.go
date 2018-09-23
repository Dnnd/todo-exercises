package models

type User struct {
	Id       string `json:"id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}

func (u User) GetMapping() string {

	return `
	{
	  "mappings": {
		"_doc": {
		  "properties": {
			"nickname": {
			  "type": "keyword"
			}
		  }
		}
	  }
	}
`
}
