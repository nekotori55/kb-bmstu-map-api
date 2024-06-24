package main

import (
	"encoding/json"
	"fmt"
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
		var output string

		var group = "ИУК4-61Б"
		var schedule = getSchedule(3, group)

		for _, daySchedule := range schedule {

			var lessons = parseDaySchedule(daySchedule, group)

			a, _ := json.MarshalIndent(lessons, "", "	")

			output += string(a[:]) + "\n"
		}

		return c.SendString(output)
	})

	log.Fatal(app.Listen(":3000"))
}

func getSchedule(year int, group string) map[string][]string {

	agent := fiber.AcquireAgent()
	agent.Request().Header.SetMethod("GET")
	agent.Request().SetRequestURI(apiPath +
		"getschedule/" +
		fmt.Sprint(year-1) + "/" +
		group + "/" +
		"")

	var err = agent.Parse()

	if err != nil {
		panic("[REQUEST ERROR] " + err.Error())
	}

	var statusCode, body, errs = agent.Bytes()

	if len(errs) > 0 {
		panic("[REQUEST ERROR] (status code: " + fmt.Sprint(statusCode) + ")")
	}

	var result map[string][]string
	err = json.Unmarshal(body, &result)

	if err != nil {
		println(err.Error())
		panic("[REQUEST ERROR] response unmarshaling error")
	}

	fiber.ReleaseAgent(agent)

	return result
}
