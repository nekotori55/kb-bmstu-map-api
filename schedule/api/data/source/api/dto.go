package api

import (
	v1 "kb-bmstu-map-api/schedule/api/domain/v1/entities"
)

type ClassDTO struct {
	Title      string
	Type       string
	Group      string
	Subgroup   int
	Building   string
	Room       string
	Professors string
	Notes      string
	Day        int
	Regularity int
	Index      int
}

func (c ClassDTO) toClassV1() v1.Class {
	return v1.Class{
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
		// Unused
		// Id:         0,
		// StartTime:  "",
		// EndTime:    "",
	}
}
