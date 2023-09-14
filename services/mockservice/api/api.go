package api

import (
	"Api_Mock/services/mockservice/data"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

func ApiRoutes(app *fiber.App, cache *cache.Cache, db *sql.DB) {
	clientStore := data.NewClientDetailsStore(db)
	requestStore := data.NewRequestStore(db)
	clientApi := ClientApi(cache, clientStore)
	requestApi := NewRequestApi(clientStore, requestStore)
	app.Post("/client", clientApi.CreateClient)
	app.Get("client/:id", clientApi.GetClientById)
	app.Post("client/:id/money-requests", requestApi.CreateMoneyRequest)
	app.Get("client/:id/money-requests", requestApi.GetMoneyRequest)
	app.Get("/money-request/:id", requestApi.GetMoneyRequestById)
	app.Put("/money-request/:id", requestApi.UpdateMoneyRequestById)
}
