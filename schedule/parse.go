package main

import (
	"github.com/dlclark/regexp2"
)

var scheduleParseExp = regexp2.MustCompile(``+
	`(?P<title>[А-Яа-я \/.,-]+?)? ?`+
	`(?P<type>упр.|лекц.|лаб.|к.р.|к.пр.)? ?`+
	`(?P<group>I|II|III|IV|IIII)? ?`+
	`(?P<location>`+
	/**/ `(?P<building>\d)[-_](?P<room>\d{3}[а-я]?(?:[\/\d]{2})?)`+
	/**/ `|к\.(?P<building>\d)`+
	/**/ `|(?P<building>УАК\d)-(?P<room>\d.\d{2})`+
	/**/ `|(?P<building>НПП "Тайфун"|ООО РИТЦ|ОКБ "МЭЛ")) ?`+
	`(?P<professors>(?:(?:[А-Я][а-я]+)+(?:, [А-Я][а-я]+)*))? ?`+
	`(?P<notes>[а-я. ]+)? ?`,
	regexp2.RE2)

func parseDaySchedule(schedule []string) []lesson {
	var romanNumbers = []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII"}
	var lessonRegularityTokens = []string{"Ч", "З", "П"}
	var weekdays = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

	var lessons []lesson

	var dayNum int
	if success, entryIndex := stringStartsWithAnyOf(weekdays, schedule[0]); success {
		dayNum = entryIndex + 1
	}

	schedule = schedule[1:] // removing weekday string

	var classNum int
	var lessonRegularity int
	for _, entry := range schedule {
		// string is a classNum
		if success, entryIndex := stringStartsWithAnyOf(romanNumbers, entry); success {
			classNum = entryIndex + 1
			continue
		}

		if slen(entry) == 1 {
			success, entryIndex := stringStartsWithAnyOf(lessonRegularityTokens, entry)
			if !success {
				panic("[PARSING ERROR] Wrong regularity token")
			}

			// BIT FLAGS
			// 'Ч'(числитель) => 1 = 01
			// 'З'(знаменатель) => 2 = 10
			// 'П'(постоянное) => 3 = 11
			lessonRegularity = entryIndex + 1
			continue
		}

		var match, match2 *regexp2.Match
		var err error

		match, err = scheduleParseExp.FindStringMatch(entry)
		if err != nil {
			panic(err.Error())
		}
		if match == nil {
			panic("[PARSING ERROR] Bad schedule entry" + entry)
		}

		var newLesson lesson = matchToLesson(match)
		newLesson.ClassNum = classNum
		newLesson.Regularity = lessonRegularity
		newLesson.Day = dayNum

		lessons = append(lessons, newLesson)

		// If the lesson string has a second part
		match2, err = scheduleParseExp.FindNextMatch(match)
		if err != nil {
			panic(err.Error())
		}
		if match2 == nil {
			continue
		}

		var newLesson2 lesson = matchToLesson(match2)
		newLesson2.ClassNum = classNum
		newLesson2.Regularity = lessonRegularity
		newLesson2.Day = dayNum

		if newLesson2.Title == "" {
			newLesson2.Title = newLesson.Title
		}

		lessons = append(lessons, newLesson2)
	}

	return lessons
}

func matchToLesson(match *regexp2.Match) lesson {
	var newLesson lesson

	newLesson.Title = match.GroupByName("title").String()
	newLesson.LessonType = match.GroupByName("type").String()
	newLesson.Group = match.GroupByName("group").String()
	newLesson.Building = match.GroupByName("building").String()
	newLesson.Room = match.GroupByName("room").String()
	newLesson.Professors = match.GroupByName("professors").String()
	newLesson.Notes = match.GroupByName("notes").String()

	return newLesson
}
