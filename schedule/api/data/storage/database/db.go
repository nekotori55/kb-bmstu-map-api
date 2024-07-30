package database

import (
	"context"
	"fmt"
	"os"

	u "kb-bmstu-map-api/schedule/api/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func setClasses(classes []ClassDTO) error {
	db.Exec(context.Background(), "TRUNCATE TABLE classes")

	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"classes"},
		[]string{
			`title`,
			`type`,
			`group`,
			`subgroup`,
			`building`,
			`room`,
			`professors`,
			`notes`,
			`day`,
			`regularity`,
			`index`,
		},
		pgx.CopyFromSlice(len(classes), func(i int) ([]any, error) {
			return []any{
				classes[i].Title,
				classes[i].Type,
				classes[i].Group,
				classes[i].Subgroup,
				classes[i].Building,
				classes[i].Room,
				classes[i].Professors,
				classes[i].Notes,
				classes[i].Day,
				classes[i].Regularity,
				classes[i].Index,
			}, nil
		}),
	)

	if copyCount != int64(len(classes)) {
		panic("not all classes were inserted")
	}
	return err
}

var db *pgx.Conn

func initDB() {
	if db != nil {
		return
	}

	var err error
	db, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func buildSearchQuery(filters map[string]string) (string, []any) {
	columns := []string{
		"id",
		"title",
		`"group"`,
		"subgroup",
		"building",
		`"type"`,
		"room",
		"professors",
		"notes",
		"regularity",
		`"day"`,
		`classes."index" as "index"`,
		"startTime",
		"endTime",
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	requestBuilder := psql.Select(columns...).From("classes").Join("time_slots ON time_slots.index = classes.index")

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
	u.Must(err)
	return query, args
}
