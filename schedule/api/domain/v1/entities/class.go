package entities

type Class struct {
	Id         int    `json:"id" db:"id" example:"208"`
	Title      string `json:"title" db:"title" example:"Моделир. процессов в техносфере"`
	Type       string `json:"type" db:"type" example:"упр."`
	Group      string `json:"group" db:"group" example:"ИУК7-21М"`
	Subgroup   int    `json:"subgroup" db:"subgroup" example:"0"`
	Building   string `json:"building" db:"building" example:"7"`
	Room       string `json:"room" db:"room" example:"302"`
	Professors string `json:"professors" db:"professors" example:"Морозенко"`
	Notes      string `json:"notes" db:"notes" example:"по расп.каф."`
	Day        int    `json:"day" db:"day" example:"1"`
	Regularity int    `json:"regularity" db:"regularity" example:"2"`
	Index      int    `json:"index" db:"index" example:"4"`
	StartTime  string `json:"startTime" db:"startTime" example:"14:15"`
	EndTime    string `json:"endTime" db:"endTime" example:"15:50"`
}
