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
	validator      *ValidateModel
	logger         *log.Logger
	config         *models.Config
}

func InitHandler(app *fiber.App, conf *models.Config, log *log.Logger) {
	handler := handler{
		validator:      NewValidateModel(),
		articleService: services.NewArticleService(conf, log),
		logger:         log,
		config:         conf,
	}

	api := app.Group("/api")
	api.Get("/get-blocks", handler.GetBlocks)
	api.Get("/search-blocks", handler.SearchBlocks)
	api.Get("/test", handler.Test)
	api.Get("/test-habr", handler.TestHabr)
}

// GetBlocks godoc
// @Summary Get all blocks
// @Description Get all blocks for the current filters with pagination
// @Tags blocks
// @Accept json
// @Produce json
// @Param type query string true "Type of content to retrieve" Enums(posts, articles, news) default(articles)
// @Param sort query string true "Sorting method" Enums(new, top) default(new)
// @Param page query int true "Page number (1-50)" default(1)
// @Param period query string false "Time period for filtering content" Enums(daily, weekly, monthly, yearly, alltime)
// @Param rate query string false "Filter by minimum rating percentage" Enums(rated0, rated10, rated25, rated50, rated100)
// @Param level query string false "Difficulty level filter" Enums(easy, medium, hard)
// @Success 200 {object} models.BlocksDTO
// @Success 400 {object} string "Bad request - invalid parameters provided"
// @Failure 500 {string} string "Internal server error"
// @Router /api/get-blocks [get]
func (h *handler) GetBlocks(c fiber.Ctx) error {
	req := &GetBlocksRequest{}
	if err := h.validator.ValidateRequest(c, req); err != nil {
		c.Response().Header.SetStatusCode(http.StatusBadRequest)
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
		h.logger.Printf("Error getting blocks: %v", err)
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Coudn't get blocks")
	}

	jsonBody, err := json.Marshal(results)
	if err != nil {
		h.logger.Printf("Error marshal blocks: %v", err)
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Coudn't marshal blocks")
	}

	c.Response().Header.SetStatusCode(http.StatusOK)
	c.Response().Header.Set("Content-Type", "application/json")
	return c.Send(jsonBody)
}

// SearchBlocks godoc
// @Summary Search all blocks by filters
// @Description Search all blocks for the current filters with pagination
// @Tags blocks
// @Accept json
// @Produce json
// @Param query query string true "Query for search" default(1)
// @Param sort query string true "Sorting method" Enums(relevance, date, rating) default(relevance)
// @Param page query int true "Page number (1-50)" default(1)
// @Success 200 {object} models.BlocksDTO
// @Success 400 {object} string "Bad request - invalid parameters provided"
// @Failure 500 {string} string "Internal server error"
// @Router /api/search-blocks [get]
func (h *handler) SearchBlocks(c fiber.Ctx) error {
	req := &SearchBlocksRequest{}
	if err := h.validator.ValidateRequest(c, req); err != nil {
		c.Response().Header.SetStatusCode(http.StatusBadRequest)
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
		h.logger.Printf("Error marshal blocks: %v", err)
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

// TestHabr godoc
// @Summary Test connection
// @Description Test connection to HABR
// @Tags test
// @Success 200 {string} string
// @Failure 500 {string} string "Internal server error"
// @Router /api/test-habr [get]
func (h *handler) TestHabr(c fiber.Ctx) error {
	response, err := http.Get(h.config.BaseUrl)
	if err != nil || response.StatusCode != 200 {
		h.logger.Printf("Error connecting to HABR: %v", err)
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Lost connection to HABR")
	}

	c.Response().Header.SetStatusCode(http.StatusOK)
	return c.SendString("Test OK")
}

// Test godoc
// @Summary Test connection
// @Description Test connection to server
// @Tags test
// @Success 200 {string} string
// @Router /api/test [get]
func (h *handler) Test(c fiber.Ctx) error {
	c.Response().Header.SetStatusCode(http.StatusOK)
	return c.SendString("Test OK")
}
