package database

import (
	"context"

	v1 "kb-bmstu-map-api/schedule/api/domain/v1/entities"

	"github.com/jackc/pgx/v5"
)

type databaseStorageRepositoryForV1 struct{}

func NewDatabaseStorageRepositoryForV1() *databaseStorageRepositoryForV1 {
	initDB()
	return &databaseStorageRepositoryForV1{}
}

func (s *databaseStorageRepositoryForV1) SetClasses(classes []v1.Class) error {
	c := []ClassDTO{}
	for _, class := range classes {
		c = append(c, toClassDtoFromV1(class))
	}

	err := setClasses(c)
	return err
}

func (s *databaseStorageRepositoryForV1) GetClassesFiltered(filters map[string]string) (c []v1.Class, err error) {
	c = []v1.Class{}
	query, args := buildSearchQuery(filters)

	rows, err := db.Query(context.Background(), query, args...)

	classes, err := pgx.CollectRows(rows, pgx.RowToStructByName[ClassDTO])

	for _, class := range classes {
		c = append(c, class.toClassV1())
	}
	return c, err
}
