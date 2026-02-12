package helper

import (
	"habrexclude/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BlockInfoHelper struct{}

func NewBlockInfoHelper() *BlockInfoHelper {
	return &BlockInfoHelper{}
}

func (bih *BlockInfoHelper) GetBlockInfo(URL string) (*models.BlockInfo, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, models.ErrBadRequest
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	mockResponse := &models.HTMLResponse{
		HTMLDoc: doc,
	}

	return bih.parseArticle(mockResponse)
}

func (bih *BlockInfoHelper) parseArticle(r *models.HTMLResponse) (*models.BlockInfo, error) {

	var mainError error
	var result *models.BlockInfo

	doc := r.HTMLDoc.(*goquery.Document)
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		if mainError != nil {
			return
		}

		id := ""
		// id, ok := s.Attr("id")
		// if !ok {
		// 	mainError = models.ErrAttrNotFound
		// 	return
		// }

		date := time.Now()
		published, ok := s.Find("a.tm-article-datetime-published time").Attr("datetime")
		if ok {
			date, _ = time.Parse("2006-01-02T15:04:05Z", published)
		}

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

		result = &models.BlockInfo{
			Id:       id,
			Types:    types,
			Title:    title,
			Author:   author,
			Views:    views,
			Duration: duration,
			Level:    level,
			Date:     date,
			Tags:     tags,
			URL:      address,
		}

	})
	return result, mainError
}
