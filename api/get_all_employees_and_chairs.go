package api

import (
	"slices"
	"time"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func GetAllEmployeesAndChairs() (EmployeesAndChairs, bool, error) {
	return LazyCache("GetAllEmployeesAndChairs", time.Hour, func() (EmployeesAndChairs, error) {
		filters, err := GetFilters()
		if err != nil {
			return EmployeesAndChairs{}, err
		}

		all := EmployeesAndChairs{
			Employees: make([]KeyValue, 0),
			Chairs:    make([]KeyValue, 0),
		}

		for _, faculty := range filters.Faculties {
			employeesAndChairs, err := GetEmployeesAndChairs(faculty.Key)
			if err != nil {
				return EmployeesAndChairs{}, err
			}
			all.Employees = append(all.Employees, employeesAndChairs.Employees...)
			all.Chairs = append(all.Chairs, employeesAndChairs.Chairs...)
		}
		all.Employees = append(all.Employees, KeyValue{
			Value: "Циліндр Давид Мурадович",
			Key:   "LDQ9WWJ0W5U8",
		})
		collator := collate.New(language.Ukrainian)

		slices.SortFunc(all.Employees, func(a, b KeyValue) int {
			return collator.CompareString(a.Value, b.Value)
		})
		all.Employees = slices.CompactFunc(all.Employees, func(a, b KeyValue) bool {
			return a.Key == b.Key
		})

		return all, nil
	})
}
