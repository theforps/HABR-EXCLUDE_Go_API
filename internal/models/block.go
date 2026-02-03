package models

import "time"

type Block struct {
	Id          string    `json:"id"`
	GlobalType  string    `json:"global_type"`
	Type        string    `json:"type"`
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

type BlockInfo struct {
	
}

const (
	Post = iota
	Article
	News
	Search
)
