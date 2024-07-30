package entities

type Class struct {
	Id         int    `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Type       string `json:"type" db:"type"`
	Group      string `json:"group" db:"group"`
	Subgroup   int    `json:"subgroup" db:"subgroup"`
	Building   string `json:"building" db:"building"`
	Room       string `json:"room" db:"room"`
	Professors string `json:"professors" db:"professors"`
	Notes      string `json:"notes" db:"notes"`
	Day        int    `json:"day" db:"day"`
	Regularity int    `json:"regularity" db:"regularity"`
	Index      int    `json:"index" db:"index"`
	StartTime  string `json:"startTime" db:"startTime"`
	EndTime    string `json:"endTime" db:"endTime"`
}
