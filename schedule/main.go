package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type lesson struct {
	Id         int    `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Type       string `json:"type" db:"type"`
	Group      string `json:"group" db:"group"`
	Subgroup   int    `json:"subgroup" db:"subgroup"`
	Building   string `json:"building" db:"building"`
	Room       string `json:"room" db:"room"`
	Professors string `json:"professors" db:"professors"`
	Notes      string `json:"notes" db:"notes"`

	Day        int `json:"day" db:"day"`
	Regularity int `json:"regularity" db:"regularity"`
	Index      int `json:"index" db:"index"`

	StartTime string `json:"startTime" db:"startTime"`
	EndTime   string `json:"endTime" db:"endTime"`
}

var apiPath = "https://schedule.iuk4.ru/api/"

var db *pgx.Conn

func main() {
	app := fiber.New()

	var err error
	db, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(context.Background())

	app.Get("/schedule/update/", func(c *fiber.Ctx) error {
		coursesCount := len(getCourses())

		output := ""

		var batch pgx.Batch

		for course := 0; course < coursesCount; course++ {
			groups := Values(getGroups(course))

			for _, group := range groups {
				schedule := getSchedule(course, group)

				// optional
				// time.Sleep(100 * time.Millisecond)

				for _, daySchedule := range schedule {
					var lessons = parseDaySchedule(daySchedule, group)

					if lessons == nil {
						continue
					}

					for _, lesson := range lessons {

						batch.QueuedQueries = append(batch.QueuedQueries, &pgx.QueuedQuery{
							SQL: `INSERT INTO schedule(` +
								`title, "group",` +
								`subgroup, building, ` +
								`"type", room, ` +
								`professors, notes, ` +
								`regularity, "day", ` +
								`"index"` +
								`) VALUES (` +
								`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11` +
								`)`,
							Arguments: []any{
								lesson.Title, lesson.Group,
								lesson.Subgroup, lesson.Building,
								lesson.Type, lesson.Room,
								lesson.Professors, lesson.Notes,
								lesson.Regularity, lesson.Day,
								lesson.Index,
							},
						})
					}
				}
			}
		}

		db.Exec(context.Background(), "TRUNCATE TABLE schedule")

		results := db.SendBatch(context.Background(), &batch)
		err := results.Close()

		if err != nil {
			panic(err.Error())
		}

		return c.SendString(output)
	})

	app.Get("/schedule/get", func(c *fiber.Ctx) error {
		filters := c.Queries()

		query := `SELECT ` +
			`id, title, "group", subgroup, building, "type", room, professors, notes, regularity, "day", ` +
			`schedule."index" AS "index", startTime, endTime ` +
			`FROM schedule ` +
			`JOIN time_slots ON time_slots."index" = schedule."index" `

		if len(filters) > 0 {
			query += `WHERE `
		}

		filtersAdded := 0

		val, ok := filters["building"]
		if ok {
			query += `building = '` + val + `' `
			filtersAdded++
		}

		val, ok = filters["room"]
		if ok {
			if filtersAdded != 0 {
				query += `AND `
			}
			query += `room = '` + val + `' `
			filtersAdded++
		}

		val, ok = filters["day"]
		if ok {
			if filtersAdded != 0 {
				query += `AND `
			}
			query += `"day" = ` + val + ` `
			filtersAdded++
		}

		val, ok = filters["regularity"]
		if ok {

			if val != "3" {
				if filtersAdded != 0 {
					query += `AND `
				}
				query += `regularity = ` + val + ` `
				filtersAdded++
			}
		}

		query += `;`

		rows, err := db.Query(context.Background(), query)

		if err != nil {
			panic(err.Error())
		}

		lessons, err := pgx.CollectRows(rows, pgx.RowToStructByName[lesson])
		if err != nil {
			panic(err.Error())
		}
		if err := rows.Err(); err != nil {
			panic(err.Error())
		}

		return c.JSON(lessons)
	})

	log.Fatal(app.Listen(":3000"))
}
