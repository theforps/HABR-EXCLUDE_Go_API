package helper

import (
	"habrexclude/internal/models"

	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

type BlocksHelper struct {
	config *models.Config
}

func NewBlocksHelper(conf *models.Config) *BlocksHelper {
	return &BlocksHelper{
		config: conf,
	}
}

func (bh *BlocksHelper) GetBlocksAsync(globalType int, pageNum int) (chan *models.Block, chan error) {
	results := make(chan *models.Block)
	errCh := make(chan error, 1)

	go func() {
		defer close(results)
		defer close(errCh)

		URL := bh.buildURL(globalType, pageNum)

		geziyor.NewGeziyor(&geziyor.Options{
			StartURLs:   []string{URL},
			LogDisabled: true,
			ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
				if r.Response.StatusCode != 200 {
					errCh <- ErrBadRequest
					return
				}

				if err := parsePreviewNodes(r, results); err != nil {
					errCh <- err
					return
				}
			},
		}).Start()
	}()

	return results, errCh
}

func (bh *BlocksHelper) buildURL(globalType, pageNum int) string {
	var url string

	switch globalType {
	case models.Article:
		url = bh.config.ArticleUrl
	case models.News:
		url = bh.config.NewsUrl
	case models.Post:
		url = bh.config.PostUrl
	default:
		url = bh.config.SearchUrl
	}

	if pageNum > 1 {
		url += "/page" + strconv.Itoa(pageNum)
	}

	return url
}

func parsePreviewNodes(r *client.Response, results chan *models.Block) error {

	var mainError error

	r.HTMLDoc.Find("article.tm-articles-list__item").Each(func(i int, s *goquery.Selection) {
		if mainError != nil {
			return
		}

		article, err := parsePreviewNode(s)
		if err != nil {
			mainError = err
			return
		}
		results <- article
	})
	return mainError
}

func parsePreviewNode(s *goquery.Selection) (*models.Block, error) {
	id, ok := s.Attr("id")
	if !ok {
		return nil, ErrAttrNotFound
	}

	published, ok := s.Find("a.tm-article-datetime-published time").Attr("datetime")
	if !ok {
		return nil, ErrAttrNotFound
	}

	date, err := time.Parse("2006-01-02T15:04:05Z", published)
	if err != nil {
		return nil, ErrInvalidDate
	}

	image, _ := s.Find("img.lead-image").Attr("src")
	author := s.Find("a.tm-user-info__username").Text()
	level := s.Find("span.tm-article-complexity__label").Text()
	duration := s.Find("span.tm-article-reading-time__label").Text()
	views := s.Find("span.tm-icon-counter__value").Text()
	title := s.Find("a.tm-title__link span").Text()

	var types []string
	s.Find("div.tm-publication-label span").Each(func(i int, s *goquery.Selection) {
		typeStr := strings.TrimSpace(s.Text())
		if typeStr != "" {
			types = append(types, typeStr)
		}
	})

	var tags []string
	s.Find("a.tm-publication-hub__link").Each(func(i int, tag *goquery.Selection) {
		tagStr := strings.Replace(tag.Text(), "*", "", -1)
		tagStr = strings.TrimSpace(tagStr)
		if tagStr != "" {
			tags = append(tags, tagStr)
		}
	})

	var description strings.Builder
	s.Find("div.article-formatted-body p").Each(func(i int, s *goquery.Selection) {
		desc := strings.TrimSpace(s.Text())
		if desc != "" {
			description.WriteString(desc + "\n\n")
		}
	})
	descriptionStr := strings.TrimSpace(description.String())

	address, ok := s.Find("a.tm-title__link").Attr("href")
	if !ok {
		return nil, ErrAttrNotFound
	}

	return &models.Block{
		Id:          id,
		Types:       types,
		Title:       title,
		Author:      author,
		Views:       views,
		Duration:    duration,
		Level:       level,
		Date:        date,
		Tags:        tags,
		Image:       image,
		Description: descriptionStr,
		URL:         address,
	}, nil
}
