package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type getStudyGroupsResponse struct {
	D struct {
		StudyGroups []KeyValue `json:"studyGroups"`
	} `json:"d"`
}

func GetStudyGroups(faculty string, educationForm string, course string) ([]KeyValue, error) {
	req, _ := http.NewRequest("GET", "https://vnz.osvita.net/WidgetSchedule.asmx/GetStudyGroups", nil)

	q := req.URL.Query()
	q.Add("callback", "")
	q.Add("aVuzId", "11613")
	q.Add("aGiveStudyTimes", "false")
	q.Add("aFacultyID", "\""+faculty+"\"")
	q.Add("aEducationForm", "\""+educationForm+"\"")
	q.Add("aCourse", "\""+course+"\"")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result getStudyGroupsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	SortKeyValue(result.D.StudyGroups)

	return result.D.StudyGroups, nil
}
