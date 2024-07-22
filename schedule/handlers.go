package main

import (
	"context"
	exAPI "kb-bmstu-map-api/schedule/external_api"
	. "kb-bmstu-map-api/schedule/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func updateHandler(c *fiber.Ctx) error {
	output := updateSchedule()
	return c.SendString(output)
}

func getScheduleHandler(c *fiber.Ctx) error {
	filters := c.Queries()
	lessons := getLessons(filters)
	return c.JSON(lessons)
}

func updateSchedule() string {
	output := ""
	var batch pgx.Batch

	courses := exAPI.GetCourses()
	for courseIndex := range courses {
		groups := Values(exAPI.GetGroups(courseIndex))

		for _, group := range groups {
			schedule := exAPI.GetSchedule(courseIndex, group)

			for _, daySchedule := range schedule {
				var lessons = parseDaySchedule(daySchedule)

				for _, lesson := range lessons {
					lesson.Group = group

					batch.QueuedQueries = append(batch.QueuedQueries, &pgx.QueuedQuery{
						SQL: `INSERT INTO schedule(` +
							`title, "group", subgroup, building, ` +
							`"type", room, professors, notes, ` +
							`regularity, "day", "index"` +
							`) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
						Arguments: []any{
							lesson.Title, lesson.Group, lesson.Subgroup, lesson.Building,
							lesson.Type, lesson.Room, lesson.Professors, lesson.Notes,
							lesson.Regularity, lesson.Day, lesson.Index,
						},
					})
				}
			}
		}
	}

	db.Exec(context.Background(), "TRUNCATE TABLE schedule")

	results := db.SendBatch(context.Background(), &batch)
	err := results.Close()
	Must(err)

	return output
}

func getLessons(filters map[string]string) []lesson {
	query, args := buildSearchQuery(filters)

	rows, err := db.Query(context.Background(), query, args...)
	Must(err)

	lessons, err := pgx.CollectRows(rows, pgx.RowToStructByName[lesson])
	Must(err)
	return lessons
}
