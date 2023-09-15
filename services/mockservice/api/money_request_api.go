package api

import (
	"Api_Mock/pkg/model"
	"Api_Mock/pkg/protocol/http/response"
	"Api_Mock/services/mockservice/data"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type MoneyRequest struct {
	clientStore  data.ClientStore
	requestStore data.RequestStore
}

type RequestApi interface {
	CreateMoneyRequest(ctx *fiber.Ctx) error
	GetMoneyRequest(ctx *fiber.Ctx) error
	GetMoneyRequestById(ctx *fiber.Ctx) error
	UpdateMoneyRequestById(ctx *fiber.Ctx) error
}

func NewRequestApi(clientStore data.ClientStore, requestStore data.RequestStore) RequestApi {
	return &MoneyRequest{
		clientStore:  clientStore,
		requestStore: requestStore,
	}
}

func (m *MoneyRequest) CreateMoneyRequest(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	requesterId := ctx.Get("x-signed-on-client")

	payload := model.MoneyRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		errorResponse := response.ErrorResponse{
			ErrorCode: "Invalid-Payload",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	requesterDetails, err := m.clientStore.GetClientDetailsById(requesterId)
	if err != nil {
		errorResponse := response.ErrorResponse{
			ErrorCode: "DB-Error",
			Message:   err.Error(),
		}
		if err == sql.ErrNoRows {
			errorResponse.Message = "requester client not found"
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse)
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	requesteeDetails, err := m.clientStore.GetClientDetailsById(payload.RequesteeId)
	if err != nil {
		errorResponse := response.ErrorResponse{
			ErrorCode: "DB-Error",
			Message:   err.Error(),
		}
		if err == sql.ErrNoRows {
			errorResponse.Message = "requestee client not found"
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse)
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	moneyRequest := &data.Request{
		ID:             uuid.New().String(),
		RequestedDate:  time.Now().String(),
		RequesterId:    requesterId,
		RequesterName:  requesterDetails.Name,
		RequesteeName:  requesteeDetails.Name,
		RequesteeId:    payload.RequesteeId,
		RequestStatus:  data.PENDING,
		ExpirationDate: time.Now().Add(15 * 24 * time.Hour).String(),
		Amount:         payload.Amount,
		InvoiceNumber:  payload.InvoiceNumber, // fixme - we should generate
		Message:        payload.Message,
	}

	createdRecord, err := m.requestStore.InsertMoneyRequest(moneyRequest)
	if err != nil {
		errorResponse := response.ErrorResponse{
			ErrorCode: "DB-Error",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	return ctx.Status(fiber.StatusOK).JSON(createdRecord)

}

func (m *MoneyRequest) GetMoneyRequest(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	clientId := ctx.Get("x-signed-on-client")
	moneyRequests, err := m.requestStore.GetRequestsByClientId(clientId)
	if err != nil {
		errorResponse := response.ErrorResponse{
			ErrorCode: "Invalid-Payload",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	if len(moneyRequests) == 0 {
		errorResponse := response.ErrorResponse{
			ErrorCode: "Invalid-Payload",
			Message:   fmt.Sprintf("no requests found for the clientId %s", clientId),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	return ctx.Status(fiber.StatusOK).JSON(moneyRequests)
}

func (m *MoneyRequest) GetMoneyRequestById(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	requestId := ctx.Params("id")
	moneyRequest, err := m.requestStore.GetMoneyRequestById(requestId)
	if err != nil {
		errorResponse := response.ErrorResponse{
			ErrorCode: "Invalid-Payload",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	return ctx.Status(fiber.StatusOK).JSON(moneyRequest)
}

func (m *MoneyRequest) UpdateMoneyRequestById(ctx *fiber.Ctx) error {
	requestId := ctx.Params("id")
	payload := model.UpdateMoneyRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		errorMessage := response.ErrorResponse{
			ErrorCode: "Invalid-Payload",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorMessage)
	}
	moneyRequest, err := m.requestStore.GetMoneyRequestById(requestId)
	if err != nil {
		errorResponse := response.ErrorResponse{
			ErrorCode: "Invalid-Payload",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if moneyRequest.RequestStatus != data.PENDING {
		errorResponse := response.ErrorResponse{
			ErrorCode: "Invalid-Request",
			Message:   fmt.Sprintf("request %s has already been accepted/declined", requestId),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	action, err := payload.Action.ToString()
	if err != nil {
		errorResponse := response.ErrorResponse{
			ErrorCode: "Invalid-Payload",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	err = m.requestStore.UpdateMoneyRequestStatusById(moneyRequest.ID, action)
	if err != nil {
		errorResponse := response.ErrorResponse{
			ErrorCode: "DB-Error",
			Message:   err.Error(),
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	moneyRequest.RequestStatus = action

	return ctx.Status(fiber.StatusOK).JSON(moneyRequest)

}
