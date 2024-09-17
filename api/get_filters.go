package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type Filters struct {
	Faculties      []KeyValue `json:"faculties"`
	EducationForms []KeyValue `json:"educForms"`
	Courses        []KeyValue `json:"courses"`
}

type getFiltersResponse struct {
	D Filters `json:"d"`
}

func GetFilters() (Filters, error) {
	req, _ := http.NewRequest("GET", "https://vnz.osvita.net/WidgetSchedule.asmx/GetStudentScheduleFiltersData", nil)

	q := req.URL.Query()
	q.Add("callback", "")
	q.Add("aVuzID", "11613")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Filters{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Filters{}, err
	}

	var result getFiltersResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Filters{}, err
	}

	SortKeyValue(result.D.Courses)
	SortKeyValue(result.D.EducationForms)
	SortKeyValue(result.D.Faculties)

	return result.D, nil
}
