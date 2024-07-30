package api

import (
	v1 "kb-bmstu-map-api/schedule/api/domain/v1/entities"
	"kb-bmstu-map-api/schedule/api/domain/v1/repositories"
)

type foreignApiDataSourceForV1 struct{}

func NewForeignApiDataSourceForV1() (s repositories.ScheduleDataSource) {
	return &foreignApiDataSourceForV1{}
}

func (f foreignApiDataSourceForV1) GetClasses() (classes []v1.Class, err error) {
	data, err := getClasses()

	for _, dto := range data {
		classes = append(classes, dto.toClassV1())
	}

	return
}
