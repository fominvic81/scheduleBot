package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type EmployeesAndChairs struct {
	Employees []KeyValue `json:"employees"`
	Chairs    []KeyValue `json:"chairs"`
}

type getEmployeesAndChairsResponse struct {
	D EmployeesAndChairs `json:"d"`
}

func GetEmployeesAndChairs(faculty string) (*EmployeesAndChairs, error) {
	req, _ := http.NewRequest("GET", "https://vnz.osvita.net/WidgetSchedule.asmx/GetEmployeeChairs", nil)

	q := req.URL.Query()
	q.Add("callback", "")
	q.Add("aVuzID", "11613")
	q.Add("aFacultyID", "\""+faculty+"\"")
	q.Add("aGiveStudyTimes", "false")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result getEmployeesAndChairsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result.D, nil
}
