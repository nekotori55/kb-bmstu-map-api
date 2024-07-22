package main

import (
	"context"
	"fmt"
	"os"

	. "kb-bmstu-map-api/schedule/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func buildSearchQuery(filters map[string]string) (string, []any) {
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
	Must(err)
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
