package helper

import (
	"encoding/xml"
	"habrexclude/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SearchRssHelper struct{}

func NewSearchRssHelper() *SearchRssHelper {
	return &SearchRssHelper{}
}

func (sh *SearchRssHelper) GetSearchResults(URL string) (chan *models.Block, chan error) {
	results := make(chan *models.Block)
	errCh := make(chan error, 1)

	go func() {
		defer close(results)
		defer close(errCh)

		if err := sh.processSearch(URL, results); err != nil {
			errCh <- err
			return
		}
	}()

	return results, errCh
}

func (sh *SearchRssHelper) processSearch(URL string, results chan *models.Block) error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ErrBadRequest
	}

	var rss models.RSS
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&rss); err != nil {
		return err
	}

	for _, item := range rss.Channel.Items {
		block := sh.parseRSSItem(item)
		results <- block
	}

	return nil
}

func (sh *SearchRssHelper) parseRSSItem(item models.Item) *models.Block {
	pubDate, err := time.Parse(time.RFC1123, item.PubDate)
	if err != nil {

		pubDate, err = time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			pubDate = time.Now()
		}
	}

	description := sh.extractTextFromHTML(item.Description)

	id := item.GUID
	if id == "" {
		id = item.Link
	}

	var types []string
	var tags []string
	for i, category := range item.Categories {
		if i < len(item.Categories)-3 {
			types = append(types, category)
		} else {
			tags = append(tags, category)
		}
	}

	image := sh.extractImageFromHTML(item.Description)

	return &models.Block{
		Id:          id,
		Types:       types,
		Title:       item.Title,
		Author:      item.Creator,
		Views:       "",
		Duration:    "",
		Level:       "",
		Date:        pubDate,
		Tags:        tags,
		Image:       image,
		Description: description,
		URL:         sh.cleanURL(item.Link),
	}
}

func (sh *SearchRssHelper) extractTextFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}

	doc.Find("script, style").Remove()

	text := doc.Text()

	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n\n", "\n")
	text = strings.ReplaceAll(text, "\n\n", "\n")

	return text
}

func (sh *SearchRssHelper) extractImageFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	if img := doc.Find("img").First(); img.Length() > 0 {
		src, _ := img.Attr("src")
		return src
	}

	return ""
}

func (sh *SearchRssHelper) cleanURL(url string) string {
	if idx := strings.Index(url, "?"); idx != -1 {
		return url[:idx]
	}
	return url
}
