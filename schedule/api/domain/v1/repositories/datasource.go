package repositories

import (
	. "kb-bmstu-map-api/schedule/api/domain/v1/entities"
)

type ScheduleDataSource interface {
	GetClasses() ([]Class, error)
}
