package services

import (
	"fmt"
	"habrexclude/internal/models"
	"habrexclude/internal/parser"
	"log"
)

type BlocksService struct {
	blockFetcher *parser.BlocksFetcher
	logger       *log.Logger
	config       *models.Config
}

func NewArticleService(conf *models.Config, log *log.Logger) *BlocksService {
	return &BlocksService{
		blockFetcher: parser.NewBlocksFetcher(),
		logger:       log,
		config:       conf,
	}
}

func (as *BlocksService) GetAll(filter *models.BlocksFilter) (*models.BlocksDTO, error) {
	URL := as.BuildURL(filter)

	var blocks []*models.Block
	var err error

	if filter.Query != "" {
		blocks, err = as.blockFetcher.Search(URL)
	} else {
		blocks, err = as.blockFetcher.GetAll(URL)
	}

	if err != nil {
		as.logger.Println(err)
		return nil, err
	}

	result := &models.BlocksDTO{
		Filter:      filter,
		Content:     blocks,
		CountBlocks: len(blocks),
	}

	return result, nil
}

func (bh *BlocksService) BuildURL(filter *models.BlocksFilter) string {
	var url string

	if filter.Query != "" {
		url = fmt.Sprintf("%s?q=%s&order_by=%s&target_type=posts&hl=ru&fl=ru&fl=ru",
			bh.config.SearchUrl, filter.Query, filter.Sort)
		

	} else {
		switch filter.Type {
		// Обычный запрос контента
		case models.ContentTypeArticle:
			url = bh.config.ArticleUrl
		case models.ContentTypeNews:
			url = bh.config.NewsUrl
		case models.ContentTypePost:
			url = bh.config.PostUrl
		}

		// Добавление сортировки и фильтров
		if filter.Sort == models.SortTop {
			url += models.SortTop + "/" + filter.Period + "/"
		} else if filter.Rate != models.ViewsAll {
			url += filter.Rate + "/"
		}

		// Уровень сложности только для статей
		if filter.Type == models.ContentTypeArticle && filter.Level != models.LevelAll {
			url += filter.Level + "/"
		}

		// Пагинация
		if filter.Page != "1" {
			url += "page" + filter.Page
		}
	}

	return url
}
