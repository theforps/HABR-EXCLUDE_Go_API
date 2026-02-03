package parser

import (
	"habrexclude/internal/models"

	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

type ParseHelper struct {
	config *models.Config
}

func NewParseHelper(conf *models.Config) *ParseHelper {
	return &ParseHelper{
		config: conf,
	}
}

// func (ph *ParseHelper) GetArticleAsync(id int) {

// }

func (ph *ParseHelper) GetArticlesAsync(globalType int, pageNum int) (chan *models.Block, error) {
	results := make(chan *models.Block)
	errChan := make(chan error, 1)

	go func() {
		defer close(results)
		defer close(errChan)

		var page, url string 
		if pageNum > 1 {
			page = "/page" + strconv.Itoa(pageNum)
		}

		switch globalType {
			case models.Article:
				url = ph.config.ArticleUrl
			case models.News:
				url = ph.config.NewsUrl
			case models.Post:
				url = ph.config.PostUrl
			default:
				url = ph.config.SearchUrl
		}

		url += page

		geziyor.NewGeziyor(&geziyor.Options{
			StartURLs:   []string{url},
			LogDisabled: true,
			ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
				if err := parsePreviewNodes(r, results, globalType); err != nil {
					errChan <- err
				}
			},
		}).Start()
	}()

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return results, nil
}

func parsePreviewNodes(r *client.Response, results chan *models.Block, globalType int) error {

	var mainError error

	r.HTMLDoc.Find("article.tm-articles-list__item").Each(func(i int, s *goquery.Selection) {
		if mainError != nil {
			return 
		}

		article, err := parsePreviewNode(s, globalType)
		if err != nil {
			mainError = err
			return
		}
		results <- article
	})
	return mainError
}

func parsePreviewNode(s *goquery.Selection, globalType int) (*models.Block, error) {
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
	postType := s.Find("div.tm-article-labels__container").Text()
	views := s.Find("span.tm-icon-counter__value").Text()
	title := s.Find("a.tm-title__link span").Text()
	tagsSelector := s.Find("a.tm-publication-hub__link")

	var tags []string
	tagsSelector.Each(func(i int, tag *goquery.Selection) {
		tagStr := strings.TrimSpace(tag.Text())
		tagStr = strings.Replace(tagStr, "*", "", -1)
		if tagStr != "" {
			tags = append(tags, tagStr)
		}
	})

	descriptionSelector := s.Find("div.article-formatted-body p")
	var description strings.Builder
	descriptionSelector.Each(func(i int, s *goquery.Selection) {
		desc := strings.TrimSpace(s.Text())
		if desc != "" {
			description.WriteString(desc)
			description.WriteString("\n\n")
		}
	})
	descriptionStr := strings.TrimSpace(description.String())

	var globalTypeStr string
	switch globalType {
	case models.Article:
		globalTypeStr = "article"
	case models.News:
		globalTypeStr = "news"
	case models.Post:
		globalTypeStr = "post"
	default:
		globalTypeStr = "search"
	}

	return &models.Block{
		Id:          id,
		Type:        postType,
		GlobalType: globalTypeStr,
		Title:       title,
		Author:      author,
		Views:       views,
		Duration:    duration,
		Level:       level,
		Date:        date,
		Tags:        tags,
		Image:       image,
		Description: descriptionStr,
	}, nil
}
