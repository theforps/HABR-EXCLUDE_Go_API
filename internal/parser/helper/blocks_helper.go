package helper

import (
	"habrexclude/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BlocksHelper struct{}

func NewBlocksHelper() *BlocksHelper {
	return &BlocksHelper{}
}

func (bh *BlocksHelper) GetBlocksAsync(URL string) (chan *models.Block, chan error) {
	results := make(chan *models.Block)
	errCh := make(chan error, 1)

	go func() {
		defer close(results)
		defer close(errCh)

		if err := bh.processHTMLParsing(URL, results); err != nil {
			errCh <- err
			return
		}
	}()

	return results, errCh
}

func parsePreviewNodes(r *models.HTMLResponse, results chan *models.Block) error {

	var mainError error

	doc := r.HTMLDoc.(*goquery.Document)
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
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
	address, _ := s.Find("a.tm-title__link").Attr("href")

	var types []string
	s.Find("div.tm-publication-label a").Each(func(i int, s *goquery.Selection) {
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

func (bh *BlocksHelper) processHTMLParsing(URL string, results chan *models.Block) error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ErrBadRequest
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	mockResponse := &models.HTMLResponse{
		HTMLDoc: doc,
	}

	return parsePreviewNodes(mockResponse, results)
}
