package usecases

import (
	. "kb-bmstu-map-api/schedule/api/domain/v1/entities"
	. "kb-bmstu-map-api/schedule/api/domain/v1/repositories"
)

type ScheduleUsecase interface {
	UpdateSchedule() error
	GetClassesByFilters(filters map[string]string) ([]Class, error)
}

type scheduleUsecase struct {
	storageRepository ScheduleStorageRepository
	sourceRepository  ScheduleDataSource
}

func NewScheduleUsecase(storageRepository ScheduleStorageRepository, sourceRepository ScheduleDataSource) ScheduleUsecase {
	return &scheduleUsecase{storageRepository: storageRepository, sourceRepository: sourceRepository}
}

func (s *scheduleUsecase) UpdateSchedule() error {
	classes, err := s.sourceRepository.GetClasses()
	if err != nil {
		return err
	}

	err = s.storageRepository.SetClasses(classes)
	if err != nil {
		return err
	}

	return nil
}

func (s *scheduleUsecase) GetClassesByFilters(filters map[string]string) ([]Class, error) {
	classes, err := s.storageRepository.GetClassesFiltered(filters)
	return classes, err
}
