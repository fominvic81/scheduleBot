package api

import (
	"errors"
	"slices"
	"sync"
	"time"
)

func GetScheduleGroup(schedule []Day, startDate time.Time, endDate time.Time) error {
	employees := make([]KeyValue, 0, 6)

	for _, day := range schedule {
		for _, class := range day.Classes {
			employee, err := GetEmployeeByName(class.Employee)
			if err != nil {
				return err
			}
			if !slices.Contains(employees, employee) {
				employees = append(employees, employee)
			}
		}
	}

	scheduleByEmployee := make(map[string][]Day)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	var goerr error = nil

	// TODO: If error stop immediately
	for _, employee := range employees {
		wg.Add(1)

		go func() {
			schedule, err := GetEmployeeSchedule(employee, startDate, endDate)

			mu.Lock()
			if err != nil {
				goerr = err
			}
			scheduleByEmployee[employee.Value] = schedule
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()

	if goerr != nil {
		return goerr
	}

	for i := range schedule {
		for j := range schedule[i].Classes {
			class := schedule[i].Classes[j]
			employee := class.Employee
			employeeSchedule, ok := scheduleByEmployee[employee]
			if !ok {
				return errors.New("failed to get schedule by employee name: " + employee)
			}
			for _, day := range employeeSchedule {
				for _, employeeClass := range day.Classes {
					if employeeClass.FullDate == class.FullDate && employeeClass.Begin == class.Begin {
						schedule[i].Classes[j].Groups = employeeClass.Groups
					}
				}
			}
		}
	}

	return nil
}
