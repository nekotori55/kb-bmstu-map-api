package external_api

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	. "kb-bmstu-map-api/schedule/utils"
)

var apiPath = "https://schedule.iuk4.ru/api/"

func GetSchedule(course int, group string) map[string][]string {
	statusCode, body := GetJSONFromURI(apiPath + "/getschedule/" + fmt.Sprint(course) + "/" + group)

	if statusCode != 200 {
		panic("func getSchedule(): status code -  " + fmt.Sprint(statusCode))
	}

	var schedule map[string][]string

	err := json.Unmarshal(body, &schedule)

	Must(err)

	return schedule
}

func GetJSONFromURI(url string) (int, []byte) {
	agent := fiber.Get(url)

	err := agent.Parse()

	Must(err)

	statusCode, body, errs := agent.Bytes()
	fiber.ReleaseAgent(agent)

	if len(errs) > 0 {
		Must(errs[0])
	}

	return statusCode, body
}

func GetCourses() map[int]string {
	statusCode, body := GetJSONFromURI(apiPath + "/getcourses")
	if statusCode != 200 {
		panic("func getCourses(): status code -  " + fmt.Sprint(statusCode))
	}

	var courses map[int]string

	err := json.Unmarshal(body, &courses)

	Must(err)

	return courses
}

func GetGroups(course int) map[string]string {
	statusCode, body := GetJSONFromURI(apiPath + "/getgroups/" + fmt.Sprint(course))

	if statusCode != 200 {
		panic("func getGroups(): status code -  " + fmt.Sprint(statusCode))
	}

	var groups map[string]string

	err := json.Unmarshal(body, &groups)

	Must(err)

	return groups
}
