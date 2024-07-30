package main

import (
	"log"

	datasource "kb-bmstu-map-api/schedule/api/data/source/api"
	storage "kb-bmstu-map-api/schedule/api/data/storage/database"
	routes_v1 "kb-bmstu-map-api/schedule/api/domain/v1/routes"
	usecases_v1 "kb-bmstu-map-api/schedule/api/domain/v1/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	storageRepository := storage.NewDatabaseStorageRepositoryForV1()
	sourceRepository := datasource.NewForeignApiDataSourceForV1()
	scheduleUsecaseV1 := usecases_v1.NewScheduleUsecase(
		storageRepository,
		sourceRepository,
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello fiber versioning boilerplate")
	})

	// API Group Versions
	v1 := app.Group("/v1")
	routes_v1.ScheduleRouter(v1, scheduleUsecaseV1)

	log.Fatal(app.Listen(":3000"))
}
