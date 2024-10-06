package api

import (
	"slices"
	"strings"
	"time"
)

type Class struct {
	StudyTime  string `json:"study_time"`
	Begin      string `json:"study_time_begin"`
	End        string `json:"study_time_end"`
	Discipline string `json:"discipline"`
	Type       string `json:"study_type"`
	Cabinet    string `json:"cabinet"`
	WeekDay    string `json:"week_day"`
	FullDate   string `json:"full_date"`
	Employee   string `json:"employee"`
	Groups     string `json:"study_group"`
	Subgroup   string `json:"study_subgroup"`
}

type Day struct {
	MessageIds []int
	WeekDay    string
	Date       string
	Classes    []Class
}

func groupByDays(classes []Class) []Day {
	classByDateTime := make(map[string]Class)
	groupsByDateTime := make(map[string][]string)

	for _, class := range classes {
		dateTime := class.FullDate + "|" + class.Begin + "|" + class.Subgroup
		classByDateTime[dateTime] = class
		if class.Groups != "" {
			groupsByDateTime[dateTime] = append(groupsByDateTime[dateTime], class.Groups)
		}
	}

	for _, groups := range groupsByDateTime {
		slices.Sort(groups)
	}

	daysByDate := make(map[string]*Day)

	for dateTime, class := range classByDateTime {
		class.Groups = strings.Join(groupsByDateTime[dateTime], ", ")
		date := class.FullDate
		day, ok := daysByDate[date]
		if !ok {
			day = &Day{
				WeekDay: class.WeekDay,
				Date:    date,
				Classes: make([]Class, 0),
			}
			daysByDate[date] = day
		}

		day.Classes = append(day.Classes, class)
	}

	days := make([]Day, 0, len(daysByDate))
	for _, day := range daysByDate {
		slices.SortFunc(day.Classes, func(a Class, b Class) int {
			return strings.Compare(a.StudyTime, b.StudyTime)
		})
		days = append(days, *day)
	}

	slices.SortFunc(days, func(a Day, b Day) int {
		at, _ := time.Parse("02.01.2006", a.Date)
		bt, _ := time.Parse("02.01.2006", b.Date)
		return at.Compare(bt)
	})

	return days
}
