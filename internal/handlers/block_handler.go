package handlers

import (
	"encoding/json"
	"habrexclude/internal/models"
	"habrexclude/internal/services"

	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type BlockHandler struct {
	blockService *services.BlocksService
	validator    *ValidateModel
	logger       *log.Logger
}

func InitBlockHandler(app *fiber.App, config *models.Config, log *log.Logger) {
	blockHandler := BlockHandler{
		validator:    NewValidateModel(),
		blockService: services.NewBlockService(config, log),
		logger:       log,
	}

	api := app.Group("/api")
	api.Get("/get-blocks", blockHandler.GetBlocks)
	api.Get("/get-block-info", blockHandler.GetBlockInfo)
	api.Get("/search-blocks", blockHandler.SearchBlocks)
}

// GetBlockInfo godoc
// @Summary Get block info
// @Description Get block info by URL
// @Tags blocks
// @Accept json
// @Produce json
// @Param block_url query string true "URL of block"
// @Success 200 {string} string
// @Success 400 {object} string "Bad request - invalid parameters provided"
// @Failure 500 {string} string "Internal server error"
// @Router /api/get-block-info [get]
func (bh *BlockHandler) GetBlockInfo(c fiber.Ctx) error {
	url := c.Query("block_url", "")
	if url == "" {
		c.Response().Header.SetStatusCode(http.StatusBadRequest)
		return c.SendString("Block URL is empty")
	}

	block, err := bh.blockService.Get(url)
	if err != nil {
		bh.logger.Printf("Error getting block info: %v", err)
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Coudn't get block info")
	}

	jsonBody, err := json.Marshal(block)
	if err != nil {
		bh.logger.Printf("Error marshal block info: %v", err)
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Coudn't marshal block info")
	}

	c.Response().Header.SetStatusCode(http.StatusOK)
	c.Response().Header.Set("Content-Type", "application/json")
	return c.Send(jsonBody)
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
func (bh *BlockHandler) GetBlocks(c fiber.Ctx) error {
	req := &GetBlocksRequest{}
	if err := bh.validator.ValidateRequest(c, req); err != nil {
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

	results, err := bh.blockService.GetAll(filter)
	if err != nil {
		bh.logger.Printf("Error getting blocks: %v", err)
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Coudn't get blocks")
	}

	jsonBody, err := json.Marshal(results)
	if err != nil {
		bh.logger.Printf("Error marshal blocks: %v", err)
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
func (h *BlockHandler) SearchBlocks(c fiber.Ctx) error {
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

	results, err := h.blockService.GetAll(filter)
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
