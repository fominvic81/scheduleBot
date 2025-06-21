package api

import (
	"errors"
	"slices"
	"sync"
	"time"
)

func GetScheduleGroups(schedule []Day, startDate time.Time, endDate time.Time) error {
	employees := make([]KeyValue, 0, 6)

	var errs error = nil

	for _, day := range schedule {
		for _, class := range day.Classes {
			if len(class.Employee) == 0 {
				continue
			}
			employee, success, err := GetEmployeeByName(class.Employee)
			if err != nil {
				errs = errors.Join(errs, err)
			}
			if success {
				if !slices.Contains(employees, employee) {
					employees = append(employees, employee)
				}
			}
		}
	}

	scheduleByEmployee := make(map[string][]Day)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for _, employee := range employees {
		wg.Add(1)

		go func() {
			defer wg.Done()

			schedule, err := GetEmployeeSchedule(employee, startDate, endDate)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = errors.Join(errs, err)
			}
			scheduleByEmployee[employee.Value] = schedule
		}()
	}
	wg.Wait()

	for i := range schedule {
		for j := range schedule[i].Classes {
			class := &schedule[i].Classes[j]
			employee := class.Employee
			if len(employee) == 0 {
				class.Groups = "Не знайдено"
				continue
			}
			employeeSchedule, ok := scheduleByEmployee[employee]
			if !ok {
				errs = errors.Join(errors.New("failed to get schedule by employee name: " + employee))
				continue
			}
			for _, day := range employeeSchedule {
				for _, employeeClass := range day.Classes {
					if employeeClass.FullDate == class.FullDate && employeeClass.Begin == class.Begin {
						class.Groups = employeeClass.Groups
					}
				}
			}
		}
	}

	return errs
}
