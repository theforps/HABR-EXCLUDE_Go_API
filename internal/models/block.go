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
	Filter      *BlocksFilter `json:"filters"`
	CountBlocks int           `json:"count_block"`
	Content     []*Block      `json:"blocks"`
}

type BlocksFilter struct {
	Sort   string
	Query  string
	Period string
	Rate   string
	Level  string
	Page   string
	Type   string
}
