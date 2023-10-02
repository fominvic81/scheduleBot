

export interface UserI {
    id: number;
    messages: number;
    firstname: string;
    lastname?: string;
    username?: string;
    faculty: string;
    educationForm: string;
    course: string;
    studyGroup: string;
}

export interface EmployeeCacheI {
    id: number;
    date: Date;
    employees: { name: string, key: string }[];
}