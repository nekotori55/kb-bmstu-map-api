package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func getSchedule(course int, group string) map[string][]string {
	statusCode, body := getJSONFromURI(apiPath + "/getschedule/" + fmt.Sprint(course) + "/" + group)

	if statusCode != 200 {
		panic("func getSchedule(): status code -  " + fmt.Sprint(statusCode))
	}

	var schedule map[string][]string

	err := json.Unmarshal(body, &schedule)

	if err != nil {
		panic(err.Error())
	}

	return schedule
}

func getJSONFromURI(url string) (int, []byte) {
	agent := fiber.Get(url)

	err := agent.Parse()

	if err != nil {
		panic(err.Error())
	}

	statusCode, body, errs := agent.Bytes()
	fiber.ReleaseAgent(agent)

	if len(errs) > 0 {
		panic(errs[0].Error())
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

	if err != nil {
		panic(err.Error())
	}

	return courses
}

func getGroups(course int) map[string]string {
	statusCode, body := getJSONFromURI(apiPath + "/getgroups/" + fmt.Sprint(course))

	if statusCode != 200 {
		panic("func getGroups(): status code -  " + fmt.Sprint(statusCode))
	}

	var groups map[string]string

	err := json.Unmarshal(body, &groups)

	if err != nil {
		panic(err.Error())
	}

	return groups
}
