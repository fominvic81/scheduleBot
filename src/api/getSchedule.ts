import axios from 'axios';
import { getAllEmployees } from './getAllEmployees';
import { getEmployeeSchedule } from './getEmployeeSchedule';

const compareTime = (timeA: string, timeB: string): boolean => {
    const [hoursA, minutesA] = timeA.split(':').map(parseInt);
    const [hoursB, minutesB] = timeB.split(':').map(parseInt);
    
    if (hoursA > hoursB) return true;
    if (hoursA < hoursB) return false;
    if (minutesA > minutesB) return true;
    return false;
}

interface ScheduleRow {
    __type: string;
    study_time: string;
    study_time_begin: string;
    study_time_end: string;
    week_day: string;
    full_date: string;
    discipline: string;
    study_type: string;
    cabinet: string;
    employee: string;
    study_subgroup: null;
}
export interface ScheduleClass {
    class: string;
    begin: string;
    end: string;
    descipline: string;
    type: string;
    cabinet: string;
    employee: string;
    groups: Promise<string[]>;
}

export interface ScheduleDay {
    weekday: string;
    date: string;
    classes: ScheduleClass[];
}

export const getSchedule = async (studyGroupId: string, startDate: Date, endDate: Date, findGroups: boolean): Promise<ScheduleDay[]> => {

    const startString = startDate.toLocaleDateString('en-GB').replace(/\//g, '.');
    const endString = endDate.toLocaleDateString('en-GB').replace(/\//g, '.');

    const response = await axios(`https://vnz.osvita.net/WidgetSchedule.asmx/GetScheduleDataX`, {
        params: {
            callback: '',
            aVuzID: 11613,
            aStudyGroupID: `"${studyGroupId}"`,
            aStartDate: `"${startString}"`,
            aEndDate: `"${endString}"`,
            aStudyTypeID: 'null',
        },
    });
    const data = response.data.d as ScheduleRow[];
    if (!data) throw new Error('Failed to get schedule');
    const employees = new Map((await getAllEmployees()).map((value) => ([value.name, value.id])));

    const dayByDate: Map<string, ScheduleDay> = new Map();
 
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
    
        
        const groupsPromise = new Promise<string[]>(async (resolve) => {
            const groups: string[] = [];
            if (findGroups) {
                const employeeId = employees.get(row.employee);
        
                if (employeeId) {
                    const employeeDay = (await getEmployeeSchedule(employeeId, day.date, day.date)).at(0);
                    if (employeeDay) {
                        const collidingClasses = employeeDay.classes.filter((value) => value.class == row.study_time);
                        groups.push(...collidingClasses.map((value) => value.studyGroup));
                    } else {
                        console.error('Could not get employee day');
                    }
                } else {
                    console.error('Could not find employee id by name');
                }
                groups.sort();
            }
            resolve(groups);
        });
    
        day.classes.push({
            class: row.study_time,
            begin: row.study_time_begin,
            end: row.study_time_end,
            descipline: row.discipline,
            type: row.study_type,
            cabinet: row.cabinet,
            employee: row.employee,
            groups: groupsPromise,
        });
    }

    const days = [...dayByDate.values()].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
    days.forEach((day) => day.classes.sort((a, b) => compareTime(a.begin, b.begin) ? 1 : -1));

    return days;
}
