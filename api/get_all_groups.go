package api

import "time"

type Group struct {
	Faculty       KeyValue
	Course        KeyValue
	EducationForm KeyValue
	Group         KeyValue
}

func GetAllStudyGroups() ([]Group, bool, error) {
	return LazyCache("GetAllGroups", time.Hour, func() ([]Group, error) {
		filters, err := GetFilters()

		if err != nil {
			return []Group{}, err
		}

		allGroups := []Group{}

		for _, faculty := range filters.Faculties {
			for _, course := range filters.Courses {
				for _, form := range filters.EducationForms {
					groups, err := GetStudyGroups(faculty.Key, form.Key, course.Key)
					if err != nil {
						return []Group{}, err
					}

					for _, group := range groups {
						allGroups = append(allGroups, Group{
							Faculty:       faculty,
							Course:        course,
							EducationForm: form,
							Group:         group,
						})
					}
				}
			}
		}

		return allGroups, nil
	})
}
