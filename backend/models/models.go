package models

type Article struct {
	NoteId       string      `json:"noteId"`
	Title        string      `json:"title"`
	DateModified string      `json:"dateModified"`
	Type         string      `json:"type"`
	Mime         string      `json:"mime"`
	Summary      string      `json:"summary"`
	Content      string      `json:"content"`
	Attributes   []Attribute `json:"attributes"`
}

type Attribute struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type APIResponse struct {
	Results []Article `json:"results"`
}
