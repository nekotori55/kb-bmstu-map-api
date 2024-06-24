package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
)

type lesson struct {
	Title      string `json:"title"`
	Type       string `json:"type"`
	Group      string `json:"group"`
	Subgroup   int    `json:"subgroup"`
	Building   string `json:"building"`
	Room       string `json:"room"`
	Professors string `json:"professors"`
	Notes      string `json:"notes"`

	Day        int `json:"day"`
	Regularity int `json:"regularity"`
	Index      int `json:"index"`
}

var apiPath = "https://schedule.iuk4.ru/api/"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		coursesCount := len(getCourses())

		output := ""

		for course := 0; course < coursesCount; course++ {
			groups := Values(getGroups(course))

			for _, group := range groups {
				schedule := getSchedule(course, group)

				// optional
				// time.Sleep(100 * time.Millisecond)

				for _, daySchedule := range schedule {
					lessons := parseDaySchedule(daySchedule, group)

					stringRepresentation, _ := json.MarshalIndent(lessons, "", " ")

					output += string(stringRepresentation[:]) + "\n"
				}
			}
		}

		return c.SendString(output)
	})

	log.Fatal(app.Listen(":3000"))
}
