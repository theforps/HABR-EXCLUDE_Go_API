package handlers

import (
	"encoding/json"
	"habrexclude/internal/models"
	"habrexclude/internal/services"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	articleService *services.BlocksService
	validator      *ValidateModel
	logger         *log.Logger
}

func InitHandler(app *fiber.App, conf *models.Config, log *log.Logger) {
	log.Println("InitHandler called")

	handler := handler{
		validator:      NewValidateModel(),
		articleService: services.NewArticleService(conf, log),
		logger:         log,
	}

	app.Get("/api/get-blocks", handler.GetBlocks)
	app.Get("/api/search-blocks", handler.SearchBlocks)
	app.Get("/api/test", handler.Test)
	log.Println("Routes registered")
}

func (h *handler) GetBlocks(c *fiber.Ctx) error {
	h.logger.Println("GetBlocks called")

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
		log.Println(err)
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

func (h *handler) SearchBlocks(c *fiber.Ctx) error {
	h.logger.Println("SearchBlocks called")

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
		h.logger.Printf("Error getting blocks: %v", err)
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

func (h *handler) GetBlockInfo(c *fiber.Ctx) error {

	return c.SendStatus(http.StatusAccepted)
}

func (h *handler) Test(c *fiber.Ctx) error {
	h.logger.Println("Test endpoint called")
	return c.SendString("Test OK")
}
