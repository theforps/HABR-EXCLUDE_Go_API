package handlers

import (
	"encoding/json"
	"fmt"
	"habrexclude/internal/models"
	"habrexclude/internal/services"
	"log"
	"net/http"
	"strconv"

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
}

func (h *handler) GetBlocks(c fiber.Ctx) error {

	pageCountStr := c.Query("page", "1")
	pageCountNum, err := strconv.Atoi(pageCountStr)
	if err != nil || pageCountNum < 1 || pageCountNum > 50 {
		c.Response().Header.SetStatusCode(http.StatusBadRequest)
		return c.SendString(fmt.Sprintf("Invalid page count (1-50): %s", err.Error()))
	}

	blockType := c.Query("type", "1")
	blockTypeNum, err := strconv.Atoi(blockType)
	if err != nil || blockTypeNum < 0 || blockTypeNum > 3 {
		c.Response().Header.SetStatusCode(http.StatusBadRequest)
		return c.SendString(fmt.Sprintf("Invalid type (0-3): %s", err.Error()))
	}

	results, err := h.articleService.GetAll(blockTypeNum, pageCountNum)
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

// GetBlockInfo

// SearchBlock
