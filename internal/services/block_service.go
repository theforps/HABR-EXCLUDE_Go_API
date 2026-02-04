package services

import (
	"habrexclude/internal/models"
	"habrexclude/internal/parser"
	"log"
)

type BlocksService struct {
	blockFetcher parser.Fetcher[interface{}, *models.Block]
	logger       *log.Logger
}

func NewArticleService(conf *models.Config, log *log.Logger) *BlocksService {
	return &BlocksService{
		blockFetcher: parser.NewBlockFetcher(conf),
		logger:       log,
	}
}

func (as *BlocksService) GetAll(globalType, page int) (*models.BlocksDTO, error) {

	blocks, err := as.blockFetcher.GetAll(globalType, page)
	if err != nil {
		as.logger.Println(err)
		return nil, err
	}

	result := &models.BlocksDTO{
		Type:       getType(globalType),
		Content:    blocks,
		PageNumber: page,
		Count:      len(blocks),
	}

	return result, nil
}

func getType(num int) string {
	switch num {
	case models.Article:
		return "articles"
	case models.News:
		return "news"
	case models.Post:
		return "posts"
	default:
		return "search"
	}
}
