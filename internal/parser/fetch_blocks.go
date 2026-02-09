package parser

import (
	"habrexclude/internal/models"
	"habrexclude/internal/parser/helper"
)

type BlocksFetcher struct {
	blocksHelper    *helper.BlocksHelper
	searchRssHelper *helper.SearchRssHelper
}

func NewBlocksFetcher() *BlocksFetcher {
	return &BlocksFetcher{
		blocksHelper:    helper.NewBlocksHelper(),
		searchRssHelper: helper.NewSearchRssHelper(),
	}
}

func (af *BlocksFetcher) GetAll(URL string) ([]*models.Block, error) {
	ch, errCh := af.blocksHelper.GetBlocksAsync(URL)

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
	// use RSS
	// ch, errCh := af.searchRssHelper.GetSearchResults(URL)

	ch, errCh := af.blocksHelper.GetBlocksAsync(URL)

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
