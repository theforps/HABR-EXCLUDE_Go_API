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
}

func InitHandler(app *fiber.App, conf *models.Config, log *log.Logger) {
	log.Println("InitHandler called")

	handler := handler{
		validator:      NewValidateModel(),
		articleService: services.NewArticleService(conf, log),
		logger:         log,
	}

	api := app.Group("/api")
	api.Get("/get-blocks", handler.GetBlocks)
	api.Get("/search-blocks", handler.SearchBlocks)
	api.Get("/test", handler.Test)
	log.Println("Routes registered")
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
	h.logger.Println("GetBlocks called")

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

func (h *handler) SearchBlocks(c fiber.Ctx) error {
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

func (h *handler) GetBlockInfo(c fiber.Ctx) error {

	return c.SendStatus(http.StatusAccepted)
}

func (h *handler) Test(c fiber.Ctx) error {
	h.logger.Println("Test endpoint called")
	return c.SendString("Test OK")
}
