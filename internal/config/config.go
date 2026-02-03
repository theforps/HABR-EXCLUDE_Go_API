package config

import (
	"habrexclude/internal/models"
	
	"os"

	"github.com/anvidev/goenv"
)

func New() *models.Config {

	err := goenv.Load()
	if err != nil {
		panic(err)
	}

	return &models.Config{
		BaseUrl: os.Getenv("BASE_URL"),
		PostUrl: os.Getenv("POST_URL"),
		NewsUrl: os.Getenv("NEWS_URL"),
		ArticleUrl: os.Getenv("ARTICLE_URL"),
		SearchUrl: os.Getenv("SEARCH_URL"),
	}
}