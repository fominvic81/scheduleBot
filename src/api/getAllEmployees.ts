import { EmployeeCache } from '../db';
import { getEmployeesAndChairs } from './getEmployeesAndChairs';
import { getFaculies } from './getFaculies';

export interface Employee {
    id: string;
    name: string;
}

const hour = 1000 * 60 * 60;

const getAllEmployeesUncached = async (): Promise<Employee[]> => {

    const employees: Employee[] = [];
    const faculties = await getFaculies();

    for (const faculty of faculties) {
        const { employees: facultyEmployees } = await getEmployeesAndChairs(faculty.Key);
        employees.push(...facultyEmployees.map((value) => ({
            id: value.Key,
            name: value.Value,
        })));
    }

    return employees;
}

const updateEmployeeCache = async (): Promise<Employee[]> => {
    const employees: Employee[] = await getAllEmployeesUncached();

    EmployeeCache.set({
        date: new Date(),
        employees: employees.map((value) => ({ key: value.id, name: value.name })),
    });

    return employees;
}

let currentPromise: Promise<Employee[]> | undefined;

export const getAllEmployees = async (): Promise<Employee[]> => {
    if (currentPromise) return currentPromise;

    const promise = new Promise<Employee[]>(async (resolve) => {
        const cache = EmployeeCache.get();
    
        if (cache) {
            if (Date.now() - cache.date.getTime() > hour) updateEmployeeCache();
            resolve(cache.employees.map((value) => ({ id: value.key, name: value.name })));
            return;
        }

        resolve(await updateEmployeeCache());
    });
    
    currentPromise = promise;
    await promise;
    currentPromise = undefined;

    return promise;
}
