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

	log.Fatal(app.Listen(":3000"))
}
