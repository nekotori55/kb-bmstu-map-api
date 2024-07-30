package routes

import (
	s "kb-bmstu-map-api/schedule/api/domain/v1/usecases"

	"github.com/gofiber/fiber/v2"
)

func ScheduleRouter(app fiber.Router, scheduleUsecase s.ScheduleUsecase) {
	app.Get("/update", updateSchedule(scheduleUsecase))
	app.Get("/get", getFilteredSchedule(scheduleUsecase))
}

func updateSchedule(scheduleUsecase s.ScheduleUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := scheduleUsecase.UpdateSchedule()
		return err
	}
}

func getFilteredSchedule(scheduleUsecase s.ScheduleUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filters := c.Queries()
		classes, err := scheduleUsecase.GetClassesByFilters(filters)
		if err != nil {
			return err
		}
		return c.JSON(classes)
	}
}
