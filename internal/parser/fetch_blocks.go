package parser

import (
	"habrexclude/internal/models"
	"habrexclude/internal/parser/helper"
)

type BlocksFetcher struct {
	helper *helper.BlocksHelper
}

func NewBlocksFetcher() *BlocksFetcher {
	return &BlocksFetcher{
		helper: helper.NewBlocksHelper(),
	}
}

func (af *BlocksFetcher) GetAll(URL string) ([]*models.Block, error) {

	ch, errCh := af.helper.GetBlocksAsync(URL)

	var results []*models.Block
	for block := range ch {
		results = append(results, block)
	}

	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	default:
	}

	return results, nil
}
