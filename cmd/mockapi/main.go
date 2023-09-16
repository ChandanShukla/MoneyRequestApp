package main

import (
	db2 "Api_Mock/pkg/db"
	"Api_Mock/services/mockservice/api"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/patrickmn/go-cache"
	"os"
	"time"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	c := cache.New(5*time.Minute, 10*time.Minute)
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api.ApiRoutes(app, c, db)
	if err := db2.SetupDB(db); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("port\t", os.Getenv("APP_PORT"))

	log.Fatal(app.Listen(":8080"))

}
