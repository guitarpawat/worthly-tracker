package model

type Header struct {
	Title string    `json:"title"`
	Links []TopLink `json:"links"`
}

type Link struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

type TopLink struct {
	Name       string   `json:"name"`
	Href       string   `json:"href"`
	Highlight  bool     `json:"highlight"`
	ChildNodes []Link   `json:"childNodes"`
	PageName   []string `json:"-"`
}
