package routes

import (
	s "kb-bmstu-map-api/schedule/api/domain/v1/usecases"

	"github.com/gofiber/fiber/v2"
)

//	@title			kb-bmstu-map API
//	@version		1.0
//	@description	API для получения расписания КФ МГТУ им. Баумана по фильтрам

// @host		localhost:3000
// @BasePath	/v1/
func ScheduleRouter(app fiber.Router, scheduleUsecase s.ScheduleUsecase) {
	app.Post("/update", updateSchedule(scheduleUsecase))
	app.Get("/get", getFilteredSchedule(scheduleUsecase))
}

// updateSchedule
//
//	@Summary		Обновить расписание (с внешнего API)
//	@Description	Отправить команду для обновления базы данных расписания с внешнего API
//	@Success		200
//	@Router			/update [post]
func updateSchedule(scheduleUsecase s.ScheduleUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := scheduleUsecase.UpdateSchedule()
		return err
	}
}

// getFilteredSchedule
//
//	@Summary		Получить расписание
//	@Description	Получить отфильтрованный массив JSON объектов описывающих пары
//	@Produce		json
//	@Success		200			{array}	entities.Class
//	@Param			building	query	string	false	"Идентификатор строения (УАК? / ? / Организация)"									example("УАК2")
//	@Param			room		query	string	false	"Идентификатор кабинета (?.?? / ??? в зависимости от здания)"						example("2.15")
//	@Param			day			query	string	false	"Номер дня недели (1 - понедельник, 2 - вторник, ...)"								example("1")
//	@Param			regularity	query	int		false	"Регулярность пар (1 - числитель, 2 - знаменатель, 3 - числитель и знаменатель)"	example(2)
//
//	@Router			/get [get]
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
