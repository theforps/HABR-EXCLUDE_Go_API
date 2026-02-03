package services

import (
	"habrexclude/internal/models"
	"habrexclude/internal/parser"
	"log"
)

type BlocksService struct {
	blockFetcher parser.Fetcher[interface{}, *models.Block]
	logger *log.Logger
}

func NewArticleService(conf *models.Config , log *log.Logger) *BlocksService {
	return &BlocksService{
		blockFetcher: parser.NewBlockFetcher(conf),
		logger: log,
	}
}

func (as *BlocksService) GetAll(globalType, page int) ([]*models.Block, error) {

	result, err := as.blockFetcher.GetAll(globalType, page)
	if err != nil {
		as.logger.Println(err)
		return nil, err
	}

	return result, nil
}