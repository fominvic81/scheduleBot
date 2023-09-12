import axios from 'axios';
import { KeyValue } from '.';


export interface EmployeesAndChairs {
    __type: string;
    chairs: KeyValue[];
    employees: KeyValue[];
    studyTypes?: KeyValue[];
}
export const getEmployeesAndChairs = async (facultyId: string): Promise<EmployeesAndChairs> => {

    const response = await axios('https://vnz.osvita.net/BetaSchedule.asmx/GetEmployeeChairs', {
        params: {
            callback: '',
            aVuzID: 11613,
            aFacultyID: `"${facultyId}"`,
            aGiveStudyTimes: 'false',
        },
    });

    const data = response.data.d;
    if (!data) throw new Error('Failed to get employees');
    
    return data;
}
