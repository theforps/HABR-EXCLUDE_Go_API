package handlers

import (
	"encoding/json"
	"habrexclude/internal/models"
	"habrexclude/internal/services"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type handler struct {
	articleService *services.BlocksService
}

func InitHandler(app *fiber.App, conf *models.Config, log *log.Logger) {

	handler := handler{
		articleService: services.NewArticleService(conf, log),
	}

	api := app.Group("/api")
	api.Get("/get-blocks", handler.GetBlocks)
	api.Get("/search-blocks", handler.SearchBlocks)
}

func (h *handler) GetBlocks(c fiber.Ctx) error {

	filter := &models.BlocksFilter{
		Sort:   c.Query("sort", models.SortNew),
		Period: c.Query("period", models.PeriodDaily),
		Rate:   c.Query("rate", models.ViewsAll),
		Level:  c.Query("level", models.LevelAll),
		Page:   c.Query("page", "1"),
		Type:   c.Query("type", models.ContentTypeArticle),
	}

	results, err := h.articleService.GetAll(filter)
	if err != nil {
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Coudn't get blocks")
	}

	jsonBody, err := json.Marshal(results)
	if err != nil {
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Coudn't marshal blocks")
	}

	c.Response().Header.SetStatusCode(http.StatusOK)
	c.Response().Header.Set("Content-Type", "application/json")
	return c.Send(jsonBody)
}

func (h *handler) SearchBlocks(c fiber.Ctx) error {

	filter := &models.BlocksFilter{
		Sort:  c.Query("sort", models.SearchSortRelevance),
		Query: c.Query("query", ""),
		Page:   c.Query("page", "1"),
	}

	results, err := h.articleService.GetAll(filter)
	if err != nil {
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Coudn't get blocks")
	}

	jsonBody, err := json.Marshal(results)
	if err != nil {
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Coudn't marshal blocks")
	}

	c.Response().Header.SetStatusCode(http.StatusOK)
	c.Response().Header.Set("Content-Type", "application/json")
	return c.Send(jsonBody)
}

func (h *handler) GetBlockInfo(c fiber.Ctx) error {

	return c.SendStatus(http.StatusAccepted)
}
