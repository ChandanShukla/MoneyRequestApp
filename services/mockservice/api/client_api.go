package api

import (
	"Api_Mock/pkg/model"
	response2 "Api_Mock/pkg/protocol/http/response"
	"Api_Mock/services/mockservice/data"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type Service struct {
	cache       *cache.Cache
	clientStore data.ClientStore
}

type Api interface {
	CreateClient(ctx *fiber.Ctx) error
	GetClientById(ctx *fiber.Ctx) error
}

func ClientApi(cache *cache.Cache, userDetailsStore data.ClientStore) Api {
	return &Service{cache: cache,
		clientStore: userDetailsStore}
}

func (m *Service) CreateClient(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	payload := model.Client{}
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(400).SendString("Bad Request")
	}
	_, err := m.clientStore.GetClientDetailsByEmail(payload.EmailAddress)
	if err == nil {
		errorResponse := response2.ErrorResponse{
			ErrorCode: "Invalid-Request",
			Message:   fmt.Sprintf("a record already exist with %s email address", payload.EmailAddress),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse)

	} else if err != nil && err != sql.ErrNoRows {
		errorResponse := response2.ErrorResponse{
			ErrorCode: "Internal Error",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	clientDetail := data.ClientDetail{
		ID:           uuid.New().String(),
		EmailAddress: payload.EmailAddress,
		Name:         payload.Name,
	}

	createdUser, err := m.clientStore.InsertClientDetails(clientDetail)
	if err != nil {
		errorResponse := response2.ErrorResponse{
			ErrorCode: "db error",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	resp := response2.ClientSuccessResponse{
		Id:           createdUser.ID,
		Name:         createdUser.Name,
		EmailAddress: createdUser.EmailAddress,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)

}

func (m *Service) GetClientById(ctx *fiber.Ctx) error {

	clientId := ctx.Params("id")
	clientDetails, err := m.clientStore.GetClientDetailsById(clientId)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse := response2.ErrorResponse{
				ErrorCode: fmt.Sprintf(`db error: no record found for client id %s`, clientId),
				Message:   err.Error(),
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
		}
		errorResponse := response2.ErrorResponse{
			ErrorCode: "db error",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	return ctx.Status(200).JSON(clientDetails)
}
