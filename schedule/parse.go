package main

import (
	"errors"
	. "kb-bmstu-map-api/schedule/utils"

	"github.com/dlclark/regexp2"
)

var scheduleParseExp = regexp2.MustCompile(``+
	`(?P<title>[А-Яа-яёЁ \/.,-]+?)? ?`+
	`(?P<type>упр.|лекц.|лаб.|к.р.|к.пр.)? ?`+
	`(?P<group>I|II|III|IV|IIII)? ?`+
	`(?P<location>`+
	/**/ `(?P<building>\d)[-_](?P<room>\d{3}[а-яё]?(?:[\/\d]{2})?)`+
	/**/ `|к\.(?P<building>\d)`+
	/**/ `|(?P<building>УАК\d)-(?P<room>\d.\d{2})`+
	/**/ `|(?P<building>НПП "Тайфун"|ООО РИТЦ|ОКБ "МЭЛ")) ?`+
	`(?P<professors>(?:(?:[А-ЯЁ][а-яё]+)+(?:, [А-ЯЁ][а-яё]+)*))? ?`+
	`(?P<notes>[а-яё. ]+)? ?`,
	regexp2.RE2)

var lessonRegularityTokens = []string{"Ч", "З", "П"}

func parseDaySchedule(schedule []string) []lesson {
	var lessons []lesson

	day := schedule[0]
	schedule = schedule[1:] // removing weekday

	dayNum, err := parseDay(day)
	Must(err)

	// BIT FLAGS						 \n
	// 'Ч'(числитель) => 1 = 01 		 \n
	// 'З'(знаменатель) => 2 = 10		 \n
	// 'П'(постоянное) => 3 = 11
	var lessonRegularity int
	var timeSlot int
	for _, entry := range schedule {
		if timeSlotIndex, err := parseTimeSlot(entry); err == nil {
			timeSlot = timeSlotIndex
			continue
		}

		if StringLen(entry) == 1 {
			if lessonRegularity, err = parseRegularity(entry); err == nil {
				continue
			}
		}

		var match, match2 *regexp2.Match
		var err error

		match, err = scheduleParseExp.FindStringMatch(entry)
		Must(err)
		if match == nil {
			panic("[PARSING ERROR] Bad schedule entry" + entry)
		}

		var newLesson lesson = matchToLesson(match)
		newLesson.Index = timeSlot
		if timeSlot == 0 {
			print(entry)
		}
		newLesson.Regularity = lessonRegularity
		newLesson.Day = dayNum

		// If the lesson string has a second part
		match2, err = scheduleParseExp.FindNextMatch(match)
		Must(err)
		if match2 == nil {
			continue
		}

		var adjacentLesson lesson = matchToLesson(match2)
		if adjacentLesson.Title == "" {
			adjacentLesson.Title = newLesson.Title
		}
		if adjacentLesson.Type == "" {
			adjacentLesson.Type = newLesson.Type
		}
		adjacentLesson.Index = timeSlot
		adjacentLesson.Regularity = lessonRegularity
		adjacentLesson.Day = dayNum

		lessons = append(lessons, newLesson)
		lessons = append(lessons, adjacentLesson)
	}

	return lessons
}

func parseDay(entry string) (dayNum int, err error) {
	if tokenIndex := StringStartsWithToken(Weekdays, entry); tokenIndex != -1 {
		dayNum = tokenIndex + 1
	} else {
		err = errors.New("[PARSING ERROR] Wrong day token: " + entry)
	}
	return
}

func parseRegularity(entry string) (regularity int, err error) {
	if tokenIndex := StringStartsWithToken(lessonRegularityTokens, entry); tokenIndex != -1 {
		regularity = tokenIndex + 1
	} else {
		err = errors.New("[PARSING ERROR] Wrong regularity token: " + entry)
	}
	return
}

func parseTimeSlot(entry string) (timeSlot int, err error) {
	if tokenIndex := StringStartsWithToken(RomanNumbers, entry); tokenIndex != -1 {
		timeSlot = tokenIndex + 1
	} else {
		err = errors.New("[PARSING ERROR] Wrong timeSlot token: " + entry)
	}
	return
}

func parseSubgroup(entry string) (subgroup int) {
	if tokenIndex := StringStartsWithToken(RomanNumbers, entry); tokenIndex != -1 {
		subgroup = tokenIndex + 1
	} else {
		subgroup = 0
	}
	return
}

func matchToLesson(match *regexp2.Match) lesson {
	var lesson lesson

	subgroupString := match.GroupByName("group").String()
	lesson.Subgroup = parseSubgroup(subgroupString)

	lesson.Title = match.GroupByName("title").String()
	lesson.Type = match.GroupByName("type").String()
	lesson.Building = match.GroupByName("building").String()
	lesson.Room = match.GroupByName("room").String()
	lesson.Professors = match.GroupByName("professors").String()
	lesson.Notes = match.GroupByName("notes").String()

	return lesson
}
