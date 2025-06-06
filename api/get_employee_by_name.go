package api

import "errors"

func GetEmployeeByName(name string) (KeyValue, error) {
	employeesAndChairs, err := GetAllEmployeesAndChairs()
	if err != nil {
		return KeyValue{}, err
	}
	employees := employeesAndChairs.Employees

	for _, employee := range employees {
		if employee.Value == name {
			return employee, nil
		}
	}
	return KeyValue{}, errors.New("failed to find employee: " + name)
}
