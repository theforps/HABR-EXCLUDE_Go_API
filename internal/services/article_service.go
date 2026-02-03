package services

import (
	"habrexclude/internal/models"
	"habrexclude/internal/parser"
)

type ArticleService struct {
	blockFetcher parser.Fetcher[interface{}, *models.Block]
}

func NewArticleService(conf *models.Config) *ArticleService {
	return &ArticleService{
		blockFetcher: parser.NewBlockFetcher(conf),
	}
}