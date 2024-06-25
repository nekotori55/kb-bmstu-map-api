package main

import (
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

func parseDaySchedule(schedule []string, group string) []lesson {
	var lessonRegularityTokens = []string{"Ч", "З", "П"}

	var lessons []lesson

	var dayNum int
	if success, entryIndex := stringStartsWithAnyOf(weekdays, schedule[0]); success {
		dayNum = entryIndex + 1
	}

	schedule = schedule[1:] // removing weekday string

	var timeSlot int
	var lessonRegularity int
	for _, entry := range schedule {
		// string is a classNum
		if success, entryIndex := stringStartsWithAnyOf(romanNumbers, entry); success {
			timeSlot = entryIndex + 1
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
		newLesson.Index = timeSlot
		newLesson.Regularity = lessonRegularity
		newLesson.Day = dayNum
		newLesson.Group = group

		lessons = append(lessons, newLesson)

		// If the lesson string has a second part
		match2, err = scheduleParseExp.FindNextMatch(match)
		if err != nil {
			panic(err.Error())
		}
		if match2 == nil {
			continue
		}

		var additionalLesson lesson = matchToLesson(match2)
		additionalLesson.Index = timeSlot
		additionalLesson.Regularity = lessonRegularity
		additionalLesson.Day = dayNum
		additionalLesson.Group = group

		if additionalLesson.Title == "" {
			additionalLesson.Title = newLesson.Title
		}
		if additionalLesson.Type == "" {
			additionalLesson.Type = newLesson.Type
		}

		lessons = append(lessons, additionalLesson)
	}

	return lessons
}

func matchToLesson(match *regexp2.Match) lesson {
	var newLesson lesson

	newLesson.Title = match.GroupByName("title").String()
	newLesson.Type = match.GroupByName("type").String()
	newLesson.Group = match.GroupByName("group").String()
	newLesson.Building = match.GroupByName("building").String()
	newLesson.Room = match.GroupByName("room").String()
	newLesson.Professors = match.GroupByName("professors").String()
	newLesson.Notes = match.GroupByName("notes").String()

	return newLesson
}
