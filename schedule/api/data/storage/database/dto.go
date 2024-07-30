package database

import v1 "kb-bmstu-map-api/schedule/api/domain/v1/entities"

type ClassDTO struct {
	Id         int    `db:"id"`
	Title      string `db:"title"`
	Type       string `db:"type"`
	Group      string `db:"group"`
	Subgroup   int    `db:"subgroup"`
	Building   string `db:"building"`
	Room       string `db:"room"`
	Professors string `db:"professors"`
	Notes      string `db:"notes"`
	Day        int    `db:"day"`
	Regularity int    `db:"regularity"`
	Index      int    `db:"index"`
	StartTime  string `db:"startTime"`
	EndTime    string `db:"endTime"`
}

func (c ClassDTO) toClassV1() v1.Class {
	return v1.Class{
		Id:         c.Id,
		Title:      c.Title,
		Type:       c.Type,
		Group:      c.Group,
		Subgroup:   c.Subgroup,
		Building:   c.Building,
		Room:       c.Room,
		Professors: c.Professors,
		Notes:      c.Notes,
		Day:        c.Day,
		Regularity: c.Regularity,
		Index:      c.Index,
		StartTime:  c.StartTime,
		EndTime:    c.EndTime,
	}
}

func toClassDtoFromV1(c v1.Class) ClassDTO {
	return ClassDTO{
		Title:      c.Title,
		Type:       c.Type,
		Group:      c.Group,
		Subgroup:   c.Subgroup,
		Building:   c.Building,
		Room:       c.Room,
		Professors: c.Professors,
		Notes:      c.Notes,
		Day:        c.Day,
		Regularity: c.Regularity,
		Index:      c.Index,
	}
}
