package api

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"time"
)

type ScheduleClass struct {
	StudyTime  string
	Begin      string
	End        string
	Discipline string
	Type       string
	Cabinet    string
	Employee   string
}

type ScheduleDay struct {
	WeekDay string
	Date    string
	Classes []ScheduleClass
}

type scheduleRow struct {
	StudyTime      string `json:"study_time"`
	StudyTimeBegin string `json:"study_time_begin"`
	StudyTimeEnd   string `json:"study_time_end"`
	WeekDay        string `json:"week_day"`
	FullDate       string `json:"full_date"`
	Discipline     string `json:"discipline"`
	StudyType      string `json:"study_type"`
	Cabinet        string `json:"cabinet"`
	Employee       string `json:"employee"`
}

type scheduleReponse struct {
	D []scheduleRow `json:"d"`
}

func GetSchedule(studyGroup string, startDate time.Time, endDate time.Time) ([]ScheduleDay, error) {

	start := startDate.Format("02.01.2006")
	end := endDate.Format("02.01.2006")

	req, err := http.NewRequest("GET", "https://vnz.osvita.net/WidgetSchedule.asmx/GetScheduleDataX", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("callback", "")
	q.Add("aVuzId", "11613")
	q.Add("aStudyGroupID", "\""+studyGroup+"\"")
	q.Add("aStudyTypeID", "null")
	q.Add("aStartDate", "\""+start+"\"")
	q.Add("aEndDate", "\""+end+"\"")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result scheduleReponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	daysByDate := make(map[string]*ScheduleDay)

	for _, row := range result.D {
		date := row.FullDate
		day, ok := daysByDate[date]
		if !ok {
			day = &ScheduleDay{
				WeekDay: row.WeekDay,
				Date:    date,
				Classes: make([]ScheduleClass, 0),
			}
			daysByDate[date] = day
		}

		day.Classes = append(day.Classes, ScheduleClass{
			StudyTime:  row.StudyTime,
			Begin:      row.StudyTimeBegin,
			End:        row.StudyTimeEnd,
			Discipline: row.Discipline,
			Type:       row.StudyType,
			Cabinet:    row.Cabinet,
			Employee:   row.Employee,
		})
	}

	days := make([]ScheduleDay, 0, len(daysByDate))
	for _, day := range daysByDate {
		days = append(days, *day)
	}

	slices.SortFunc(days, func(a ScheduleDay, b ScheduleDay) int {
		at, _ := time.Parse("02.01.2006", a.Date)
		bt, _ := time.Parse("02.01.2006", b.Date)
		return at.Compare(bt)
	})

	return days, nil
}
