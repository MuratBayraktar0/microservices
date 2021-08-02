package main

import (
	"fmt"
	"os"

	"example.com/category/api"
	"example.com/category/repository"
	"example.com/category/service"
	"github.com/gofiber/fiber"
)

func main() {
	fmt.Println("category service starting...")
	SetupAppConfig()

	repository, err := repository.NewRepository(appConfig.DBUrl)
	if err != nil {
		fmt.Println("main.go => Repository Error: ", err)
	}
	service := service.NewService(repository, appConfig.Port)
	api := api.NewAPI(service)

	app := SetupService(api)

	fmt.Println("category service started at ", appConfig.Port, "...")
	app.Listen(appConfig.Port)
}

func SetupService(api *api.Api) *fiber.App {
	app := fiber.New()

	app.Get("/status", func(c *fiber.Ctx) {
		c.Status(fiber.StatusOK)
	})

	app.Post("/category", api.CategoryPostEndpoint)
	app.Get("/category/:categoryId", api.CategoryGetEndpoint)
	app.Get("/category/product/:productId", api.ProductCategoryGetEndpoint)

	return app
}

type Config struct {
	Port  int
	Host  string
	DBUrl string
}

var appConfig Config

func SetupAppConfig() {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "" {
		fmt.Println("APP_ENV is not set")
		os.Exit(1)
	}

	switch appEnv {
	case "test":
		appConfig = Config{
			Port:  8080,
			Host:  "http://localhost:8080",
			DBUrl: "mongodb://mongodb:27017",
		}
		return
	case "qa":
		appConfig = Config{
			Port:  8081,
			Host:  "http://localhost:8081",
			DBUrl: "mongodb://mongodb:27017",
		}
		return
	case "prod":
		appConfig = Config{
			Port:  8082,
			Host:  "http://localhost:8082",
			DBUrl: "mongodb://mongodb:27017",
		}
		return
	default:
		os.Exit(1)
	}

	appConfig = Config{}
}
