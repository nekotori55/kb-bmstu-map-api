package main

import (
	"context"
	"fmt"
	"log"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app.Use(cors.New())

	initDB()
	defer db.Close(context.Background())

	app.Get("/schedule/update/", func(c *fiber.Ctx) error {
		output := updateSchedule()
		return c.SendString(output)
	}) 

	app.Get("/schedule/get", func(c *fiber.Ctx) error {
		filters := c.Queries()
		lessons := getLessons(filters)
		return c.JSON(lessons)
	})

	log.Fatal(app.Listen(":3000"))
}

func updateSchedule() string {
	output := ""
	var batch pgx.Batch

	courses := getCourses()
	for courseIndex := range courses {
		groups := Values(getGroups(courseIndex))

		for _, group := range groups {
			schedule := getSchedule(courseIndex, group)

			for _, daySchedule := range schedule {
				var lessons = parseDaySchedule(daySchedule, group)

				for _, lesson := range lessons {

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

	if err != nil {
		panic(err.Error())
	}
	return output
}

func getLessons(filters map[string]string) []lesson {
	query, args := buildQuery(filters)

	println(query)

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		panic(err.Error())
	}

	lessons, err := pgx.CollectRows(rows, pgx.RowToStructByName[lesson])
	if err != nil {
		panic(err.Error())
	}
	return lessons
}

func buildQuery(filters map[string]string) (string, []any) {
	columns := []string{
		"id", "title", `"group"`, "subgroup", "building",
		`"type"`, "room", "professors", "notes",
		"regularity", `"day"`, `schedule."index" as "index"`, "startTime", "endTime",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	requestBuilder := psql.Select(columns...).From("schedule").Join("time_slots ON time_slots.index = schedule.index")

	if val, ok := filters["building"]; ok {
		requestBuilder = requestBuilder.Where("building = ?", val)
	}

	if val, ok := filters["room"]; ok {
		requestBuilder = requestBuilder.Where("room = ?", val)
	}

	if val, ok := filters["day"]; ok {
		requestBuilder = requestBuilder.Where(`"day" = ?`, val)
	}

	if val, ok := filters["regularity"]; ok {
		requestBuilder = requestBuilder.Where("(regularity = ? OR regularity = 3)", val)
	}

	query, args, err := requestBuilder.ToSql()
	if err != nil {
		panic(err.Error())
	}
	return query, args
}

func initDB() {
	var err error
	db, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}
