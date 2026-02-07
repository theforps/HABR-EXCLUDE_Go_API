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
	validator *ValidateModel
}

func InitHandler(app *fiber.App, conf *models.Config, log *log.Logger) {

	handler := handler{
		validator: NewValidateModel(),
		articleService: services.NewArticleService(conf, log),
	}

	api := app.Group("/api")
	api.Get("/get-blocks", handler.GetBlocks)
	api.Get("/search-blocks", handler.SearchBlocks)
}

func (h *handler) GetBlocks(c fiber.Ctx) error {

	req := &GetBlocksRequest{}
	if err := h.validator.ValidateRequest(c, req); err != nil {
		return err
	}

	filter := &models.BlocksFilter{
		Sort:   req.Sort,
		Period: req.Period,
		Rate:   req.Rate,
		Level:  req.Level,
		Page:   req.Page,
		Type:   req.Type,
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

	req := &SearchBlocksRequest{}
	if err := h.validator.ValidateRequest(c, req); err != nil {
		return err
	}

	filter := &models.BlocksFilter{
		Sort:  req.Sort,
		Query: req.Query,
		Page:  req.Page,
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
