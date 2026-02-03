package parser

import "habrexclude/internal/models"

type BlockFetcher struct {
	config *models.Config
	helper *ParseHelper
}

func NewBlockFetcher(conf *models.Config) *BlockFetcher {
	return &BlockFetcher{
		config: conf,
		helper: NewParseHelper(conf),
	}
}

func (af *BlockFetcher) GetById(id string) (interface{}, error) {

	test := &models.Block{}

	return test, nil
}

func (af *BlockFetcher) GetAll(globalType int, page int) ([]*models.Block, error) {

	ch, err := af.helper.GetArticlesAsync(globalType, page)
	if err != nil {
		return nil, err
	}

	var results []*models.Block
	for block := range ch {
		results = append(results, block)
	}

	return results, nil
}
