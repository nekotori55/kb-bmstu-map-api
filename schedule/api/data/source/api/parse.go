package api

import (
	u "kb-bmstu-map-api/schedule/api/utils"

	"github.com/dlclark/regexp2"
)

var scheduleParseExp = regexp2.MustCompile(``+
	`(?P<title>[А-Яа-яёЁ \/.,-]+?)? ?`+
	`(?P<type>упр.|лекц.|лаб.|к.р.|к.пр.)? ?`+
	`(?P<subgroup>I|II|III|IV|IIII)? ?`+
	`(?P<location>`+
	/**/ `(?P<building>\d)[-_](?P<room>\d{3}[а-яё]?(?:[\/\d]{2})?)`+
	/**/ `|к\.(?P<building>\d)`+
	/**/ `|(?P<building>УАК\d)-(?P<room>\d.\d{2})`+
	/**/ `|(?P<building>НПП "Тайфун"|ООО РИТЦ|ОКБ "МЭЛ")) ?`+
	`(?P<professors>(?:(?:[А-ЯЁ][а-яё]+)+(?:, [А-ЯЁ][а-яё]+)*))? ?`+
	`(?P<notes>[а-яё. ]+)? ?`,
	regexp2.RE2)

var lessonRegularityTokens = []string{"Ч", "З", "П"}
var RomanNumbers = []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII"}
var Weekdays = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

func parseDaySchedule(schedule []string) []ClassDTO {
	var classes []ClassDTO

	day := schedule[0]
	schedule = schedule[1:] // excluding day

	dayNum, success := parseDay(day)
	u.MustTrue(success, "[PARSING ERROR] Error in day parsing")

	var lessonRegularity int
	var timeSlot int
	for _, entry := range schedule {
		if timeSlotIndex, success := parseTimeSlot(entry); success {
			timeSlot = timeSlotIndex
			continue
		}

		if u.StringLen(entry) == 1 {
			if lessonRegularityIndex, success := parseRegularity(entry); success {
				lessonRegularity = lessonRegularityIndex
				continue
			}
		}

		var match, match2 *regexp2.Match
		var err error

		match, err = scheduleParseExp.FindStringMatch(entry)
		u.Must(err)
		if match == nil {
			panic("[PARSING ERROR] Bad schedule entry" + entry)
		}

		var newLesson ClassDTO = matchToLesson(match)
		newLesson.Index = timeSlot
		newLesson.Regularity = lessonRegularity
		newLesson.Day = dayNum

		classes = append(classes, newLesson)

		// If the lesson string has a second part
		match2, err = scheduleParseExp.FindNextMatch(match)
		u.Must(err)
		if match2 == nil {
			continue
		}

		var adjacentLesson ClassDTO = matchToLesson(match2)
		if adjacentLesson.Title == "" {
			adjacentLesson.Title = newLesson.Title
		}
		if adjacentLesson.Type == "" {
			adjacentLesson.Type = newLesson.Type
		}
		adjacentLesson.Index = timeSlot
		adjacentLesson.Regularity = lessonRegularity
		adjacentLesson.Day = dayNum

		classes = append(classes, adjacentLesson)
	}

	return classes
}

func parseDay(entry string) (dayNum int, success bool) {
	tokenIndex, success := u.StringStartsWithToken(Weekdays, entry)
	dayNum = tokenIndex + 1
	return
}

func parseRegularity(entry string) (regularity int, success bool) {
	tokenIndex, success := u.StringStartsWithToken(lessonRegularityTokens, entry)
	regularity = tokenIndex + 1
	return
}

func parseTimeSlot(entry string) (timeSlot int, success bool) {
	tokenIndex, success := u.StringStartsWithToken(RomanNumbers, entry)
	timeSlot = tokenIndex + 1
	return
}

func parseSubgroup(entry string) (subgroup int, success bool) {
	tokenIndex, success := u.StringStartsWithToken(RomanNumbers, entry)
	subgroup = tokenIndex + 1
	return
}

func matchToLesson(match *regexp2.Match) ClassDTO {
	var lesson ClassDTO

	subgroupString := match.GroupByName("subgroup").String()
	lesson.Subgroup, _ = parseSubgroup(subgroupString)

	lesson.Title = match.GroupByName("title").String()
	lesson.Type = match.GroupByName("type").String()
	lesson.Building = match.GroupByName("building").String()
	lesson.Room = match.GroupByName("room").String()
	lesson.Professors = match.GroupByName("professors").String()
	lesson.Notes = match.GroupByName("notes").String()

	return lesson
}
