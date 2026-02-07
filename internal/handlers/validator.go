package handlers

import (
	"habrexclude/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type GetBlocksRequest struct {
	Sort   string `query:"sort" validate:"omitempty,oneof=new top"`
	Period string `query:"period" validate:"omitempty,oneof=daily weekly monthly yearly alltime"`
	Rate   string `query:"rate" validate:"omitempty,oneof= rated0 rated10 rated25 rated50 rated100"`
	Level  string `query:"level" validate:"omitempty,oneof= easy medium hard"`
	Page   string `query:"page" validate:"omitempty,numeric,min=1,max=50"`
	Type   string `query:"type" validate:"omitempty,oneof=posts articles news"`
}

type SearchBlocksRequest struct {
	Query string `query:"query" validate:"required,min=1,max=100"`
	Sort  string `query:"sort" validate:"omitempty,oneof=relevance date rating"`
	Page  string `query:"page" validate:"omitempty,numeric,min=1"`
}

type ValidateModel struct{
	validate *validator.Validate
}

func NewValidateModel() *ValidateModel {
	return &ValidateModel{
		validate: validator.New(),
	}
}

func (v *ValidateModel) ValidateRequest(c fiber.Ctx, req interface{}) error {
    switch r := req.(type) {
    case *GetBlocksRequest:
        r.Sort = c.Query("sort", models.SortNew)
        r.Period = c.Query("period", models.PeriodDaily)
        r.Rate = c.Query("rate", models.ViewsAll)
        r.Level = c.Query("level", models.LevelAll)
        r.Page = c.Query("page", "1")
        r.Type = c.Query("type", models.ContentTypeArticle)
    case *SearchBlocksRequest:
        r.Query = c.Query("query", "")
        r.Sort = c.Query("sort", models.SearchSortRelevance)
        r.Page = c.Query("page", "1")
    default:
        return fiber.NewError(fiber.StatusBadRequest, "Unsupported request type")
    }

    if err := v.validate.Struct(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }

    return nil
}