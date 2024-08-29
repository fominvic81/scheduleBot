package api

var cached *EmployeesAndChairs

func GetAllEmployeesAndChairs() (*EmployeesAndChairs, error) {
	if cached != nil {
		return cached, nil
	}

	filters, err := GetFilters()
	if err != nil {
		return nil, err
	}

	all := EmployeesAndChairs{
		Employees: make([]KeyValue, 0),
		Chairs:    make([]KeyValue, 0),
	}

	for _, faculty := range filters.Faculties {
		employeesAndChairs, err := GetEmployeesAndChairs(faculty.Key)
		if err != nil {
			return nil, err
		}
		all.Employees = append(all.Employees, employeesAndChairs.Employees...)
		all.Chairs = append(all.Chairs, employeesAndChairs.Chairs...)
	}

	cached = &all
	return cached, nil
}
