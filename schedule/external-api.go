package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var apiPath = "https://schedule.iuk4.ru/api/"

func getSchedule(course int, group string) map[string][]string {
	statusCode, body := getJSONFromURI(apiPath + "/getschedule/" + fmt.Sprint(course) + "/" + group)

	if statusCode != 200 {
		panic("func getSchedule(): status code -  " + fmt.Sprint(statusCode))
	}

	var schedule map[string][]string

	err := json.Unmarshal(body, &schedule)

	Check(err)

	return schedule
}

func getJSONFromURI(url string) (int, []byte) {
	agent := fiber.Get(url)

	err := agent.Parse()

	Check(err)

	statusCode, body, errs := agent.Bytes()
	fiber.ReleaseAgent(agent)

	if len(errs) > 0 {
		Check(errs[0])
	}

	return statusCode, body
}

func getCourses() map[int]string {
	statusCode, body := getJSONFromURI(apiPath + "/getcourses")
	if statusCode != 200 {
		panic("func getCourses(): status code -  " + fmt.Sprint(statusCode))
	}

	var courses map[int]string

	err := json.Unmarshal(body, &courses)

	Check(err)

	return courses
}

func getGroups(course int) map[string]string {
	statusCode, body := getJSONFromURI(apiPath + "/getgroups/" + fmt.Sprint(course))

	if statusCode != 200 {
		panic("func getGroups(): status code -  " + fmt.Sprint(statusCode))
	}

	var groups map[string]string

	err := json.Unmarshal(body, &groups)

	Check(err)

	return groups
}
