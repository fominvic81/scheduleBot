package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type scheduleResponse struct {
	D []Class `json:"d"`
}

func GetSchedule(studyGroup string, startDate time.Time, endDate time.Time) ([]Day, error) {
	start := startDate.Format("02.01.2006")
	end := endDate.Format("02.01.2006")

	req, _ := http.NewRequest("GET", "https://vnz.osvita.net/WidgetSchedule.asmx/GetScheduleDataX", nil)

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

	var result scheduleResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return groupByDays(result.D), nil
}
