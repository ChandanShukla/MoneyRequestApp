package main

import (
	db2 "Api_Mock/pkg/db"
	"Api_Mock/services/mockservice/api"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/patrickmn/go-cache"
	"os"
	"time"
)

func main() {
	app := fiber.New()
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
