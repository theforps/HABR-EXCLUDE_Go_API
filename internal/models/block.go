package models

import "time"

type Block struct {
	Id          string    `json:"id"`
	Types       []string  `json:"types"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Views       string    `json:"views"`
	Duration    string    `json:"duration"`
	Level       string    `json:"level"`
	Date        time.Time `json:"published_at"`
	Tags        []string  `json:"tags"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
}

type BlocksDTO struct {
	Type       string   `json:"type"`
	PageNumber int      `json:"page"`
	Count      int      `json:"current_count"`
	Content    []*Block `json:"blocks"`
}

const (
	Post = iota
	Article
	News
	Search
)
