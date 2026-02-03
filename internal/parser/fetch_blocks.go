package parser

import (
	"habrexclude/internal/models"
	"habrexclude/internal/parser/helper"
)

type BlockFetcher struct {
	config *models.Config
	helper *helper.BlocksHelper
}

func NewBlockFetcher(conf *models.Config) *BlockFetcher {
	return &BlockFetcher{
		config: conf,
		helper: helper.NewBlocksHelper(conf),
	}
}

func (af *BlockFetcher) GetById(id string) (interface{}, error) {

	test := &models.Block{}

	return test, nil
}

func (af *BlockFetcher) GetAll(globalType int, page int) ([]*models.Block, error) {

	ch, errCh := af.helper.GetBlocksAsync(globalType, page)

	var results []*models.Block
	for block := range ch {
		results = append(results, block)
	}

	select {
	case err := <- errCh:
		if err != nil {
			return nil, err
		}
	default:
	}

	return results, nil
}
