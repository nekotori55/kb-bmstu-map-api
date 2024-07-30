package api

import (
	"encoding/json"
	"fmt"

	u "kb-bmstu-map-api/schedule/api/utils"

	"github.com/gofiber/fiber/v2"
)

const apiPath = "https://schedule.iuk4.ru/api/"

func getClasses() ([]ClassDTO, error) {
	classes := []ClassDTO{}

	courses := getCourses()
	for courseIndex := range courses {
		groups := u.Values(getGroups(courseIndex))

		for _, group := range groups {
			schedule := getSchedule(courseIndex, group)

			for _, daySchedule := range schedule {
				var dayClasses = parseDaySchedule(daySchedule)
				for _, class := range dayClasses {
					class.Group = group
					classes = append(classes, class)
				}
			}
		}
	}

	return classes, nil
}

func getSchedule(course int, group string) map[string][]string {

	statusCode, body := getJSONFromURI(apiPath + "getschedule/" + fmt.Sprint(course) + "/" + group)

	if statusCode != 200 {
		panic("func getSchedule(): status code -  " + fmt.Sprint(statusCode))
	}

	var schedule map[string][]string

	err := json.Unmarshal(body, &schedule)

	u.Must(err)

	return schedule
}

func getJSONFromURI(url string) (int, []byte) {
	agent := fiber.Get(url)
	err := agent.Parse()
	u.Must(err)

	statusCode, body, errs := agent.Bytes()

	if len(errs) > 0 {
		u.Must(errs[0], url)
	}

	return statusCode, body
}

func getCourses() map[int]string {
	statusCode, body := getJSONFromURI(apiPath + "getcourses")
	if statusCode != 200 {
		panic("func getCourses(): status code -  " + fmt.Sprint(statusCode))
	}

	var courses map[int]string

	err := json.Unmarshal(body, &courses)

	u.Must(err)

	return courses
}

func getGroups(course int) map[string]string {
	statusCode, body := getJSONFromURI(apiPath + "getgroups/" + fmt.Sprint(course))

	if statusCode != 200 {
		panic("func getGroups(): status code -  " + fmt.Sprint(statusCode))
	}

	var groups map[string]string

	err := json.Unmarshal(body, &groups)

	u.Must(err)

	return groups
}
