import axios from 'axios';

interface EmployeeScheduleRow {
    __type: string;
    study_time: string;
    study_time_begin: string;
    study_time_end: string;
    week_day: string;
    full_date: string;
    discipline: string;
    study_type: string;
    cabinet: string;
    study_group: string;
}
export interface EmployeeScheduleClass {
    class: string;
    begin: string;
    end: string;
    descipline: string;
    type: string;
    cabinet: string;
    studyGroup: string;
}
export interface EmployeeScheduleDay {
    weekday: string;
    date: string;
    classes: EmployeeScheduleClass[];
}

export const getEmployeeSchedule = async (employeeId: string, startDate: string, endDate: string): Promise<EmployeeScheduleDay[]> => {

    const response = await axios('https://vnz.osvita.net/WidgetSchedule.asmx/GetScheduleDataEmp', {
        params: {
            callback: '',
            aVuzID: 11613,
            aEmployeeID: `"${employeeId}"`,
            aStartDate: `"${startDate}"`,
            aEndDate: `"${endDate}"`,
            aStudyTypeID: 'null',
        },
    });

    const data = response.data.d as EmployeeScheduleRow[];
    if (!data) throw new Error('Failed to get employee schedule');

    const dayByDate: Map<string, EmployeeScheduleDay> = new Map();
    
    for (const row of data) {
        let day = dayByDate.get(row.full_date);
        if (!day) {
            day = {
                weekday: row.week_day,
                date: row.full_date,
                classes: [],
            };
            dayByDate.set(day.date, day);
        }

        day.classes.push({
            class: row.study_time,
            begin: row.study_time_begin,
            end: row.study_time_end,
            descipline: row.discipline,
            type: row.study_type,
            cabinet: row.cabinet,
            studyGroup: row.study_group,
        });
    }

    const days = [...dayByDate.values()].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());

    return days;
}