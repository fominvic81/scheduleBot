package api

import "errors"

func GetEmployeeByName(name string) (KeyValue, bool, error) {
	employeesAndChairs, success, err := GetAllEmployeesAndChairs()
	if !success {
		return KeyValue{}, false, err
	}
	employees := employeesAndChairs.Employees

	for _, employee := range employees {
		if employee.Value == name {
			return employee, true, err
		}
	}
	return KeyValue{}, false, errors.New("failed to find employee: " + name)
}
