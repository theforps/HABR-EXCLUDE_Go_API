package parser

import (
	"habrexclude/internal/models"
	"habrexclude/internal/parser/helper"
)

type BlocksFetcher struct {
	helper       *helper.BlocksHelper
	searchHelper *helper.SearchHelper
}

func NewBlocksFetcher() *BlocksFetcher {
	return &BlocksFetcher{
		helper:       helper.NewBlocksHelper(),
		searchHelper: helper.NewSearchHelper(),
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
	case <-make(chan struct{}, 1):
	}

	return results, nil
}

func (af *BlocksFetcher) Search(URL string) ([]*models.Block, error) {
	ch, errCh := af.searchHelper.GetSearchResults(URL)

	var results []*models.Block
	for block := range ch {
		results = append(results, block)
	}

	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-make(chan struct{}, 1):
	}

	return results, nil
}
