package parser

import (
	"habrexclude/internal/models"
	"habrexclude/internal/parser/helper"
)

type BlocksFetcher struct {
	blocksHelper    *helper.BlocksHelper
	blockInfoHelper *helper.BlockInfoHelper
	searchRssHelper *helper.SearchRssHelper
}

func NewBlocksFetcher() *BlocksFetcher {
	return &BlocksFetcher{
		blocksHelper:    helper.NewBlocksHelper(),
		searchRssHelper: helper.NewSearchRssHelper(),
		blockInfoHelper: helper.NewBlockInfoHelper(),
	}
}

func (bf *BlocksFetcher) GetBlockInfo(URL string) (*models.BlockInfo, error) {
	block, err := bf.blockInfoHelper.GetBlockInfo(URL)

	if err != nil {
		return nil, err
	}
	return block, nil
}

func (bf *BlocksFetcher) GetAll(URL string) ([]*models.Block, error) {
	ch, errCh := bf.blocksHelper.GetBlocksAsync(URL)

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

func (bf *BlocksFetcher) Search(URL string) ([]*models.Block, error) {
	// use RSS
	// ch, errCh := bf.searchRssHelper.GetSearchResults(URL)

	ch, errCh := bf.blocksHelper.GetBlocksAsync(URL)

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
