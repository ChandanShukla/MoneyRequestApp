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
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/client", clientApi.GetClientByEmail)
	v1.Post("/client", clientApi.CreateClient)
	v1.Get("client/:id", clientApi.GetClientById)
	v1.Post("/money-request", requestApi.CreateMoneyRequest)
	v1.Get("/money-request", requestApi.GetMoneyRequest)
	v1.Get("/money-request/:id", requestApi.GetMoneyRequestById)
	v1.Put("/money-request/:id", requestApi.UpdateMoneyRequestById)
}
