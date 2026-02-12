package models

import "time"

type BlockInfo struct {
	Id       string    `json:"id"`
	Types    []string  `json:"types"`
	URL      string    `json:"url"`
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Views    string    `json:"views"`
	Duration string    `json:"duration"`
	Level    string    `json:"level"`
	Date     time.Time `json:"published_at"`
	Tags     []string  `json:"tags"`
	Boxes    []*Box    `json:"boxes"`
}

type Box struct {
	Type     string
	Content  string
	ParentId string
}

const (
	Code        = "code"
	Enum        = "enum"
	BigHeader   = "h1"
	SmallHeader = "h2"
	Paragraph   = "p"
	Image       = "image"
	Quote       = "quote"
	Secret      = "secret"
)
