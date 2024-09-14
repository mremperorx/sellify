package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres"
	"github.com/mremperorx/sellify/api"
	"github.com/mremperorx/sellify/db"
	"github.com/mremperorx/sellify/middleware"
)

const EXPIRATION = 24 * time.Hour * 7

func init() {
	db.ConnectDB()
	err := db.CreateTables()
	if err != nil {
		fmt.Println(err)
		panic(`Failed to create tables`)
	}
	err = db.InesrtDataInCategory()
	if err != nil {
		fmt.Println(err)
		panic(`Failed to insert data`)
	}
}

func main() {
	postgresStorage := postgres.New(postgres.Config{
		Database: "sellify",
		Table:    "session",
		Username: "postgres",
		Password: "alnaser0",
	})
	sessionStore := session.New(session.Config{
		CookieHTTPOnly: true,
		Expiration:     EXPIRATION,
		Storage:        postgresStorage,
		CookieSameSite: "None",
	})

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET, POST",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	app.Post("/api/auth/signup", api.SignUpRequest)
	app.Post("/api/auth/login", api.LoginRequest(sessionStore))
	app.Post("/api/ads", middleware.Authentication(sessionStore), api.PostRequest())
	err := app.Listen(":8080")

	if err != nil {
		log.Fatal(err)
	}
}
