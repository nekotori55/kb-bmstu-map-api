package repositories

import (
	. "kb-bmstu-map-api/schedule/api/domain/v1/entities"
)

type ScheduleStorageRepository interface {
	SetClasses(classes []Class) error
	GetClassesFiltered(filters map[string]string) ([]Class, error)
}
