import axios from 'axios';

export interface KeyValue {
    Key: string;
    Value: string;
}

export interface FiltersData {
    __type: string;
    faculties: KeyValue[];
    educForms: KeyValue[];
    courses: KeyValue[];
}

export const getFiltersData = async () => {
    
    const response = await axios(`https://vnz.osvita.net/BetaSchedule.asmx/GetStudentScheduleFiltersData?callback=&aVuzID=11613`);
    const data = response.data.d as FiltersData;
    if (!data) throw new Error('Failed to get filters');

    return data;
}

export const getStudyGroupByFilters = async (facultyKey: string, educationFormKey: string, courseKey: string) => {
    
    const response = await axios(`https://vnz.osvita.net/BetaSchedule.asmx/GetStudyGroups`, {
        params: {
            callback: '',
            aVuzID: 11613,
            aFacultyID: `"${facultyKey}"`,
            aEducationForm: `"${educationFormKey}"`,
            aCourse: `"${courseKey}"`,     
            aGiveStudyTimes: false,
        },
    });
    const data = response.data.d?.studyGroups as KeyValue[];
    if (!data) throw new Error('Failed to get study groups');
    
    return data;

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
}
export interface ScheduleDay {
    weekday: string;
    date: string;
    classes: ScheduleClass[];
}

export const getSchedule = async (studyGroupKey: string, startDate: Date, endDate: Date) => {

    const startString = startDate.toLocaleDateString('en-GB').replace(/\//g, '.');
    const endString = endDate.toLocaleDateString('en-GB').replace(/\//g, '.');

    const response = await axios(`https://vnz.osvita.net/BetaSchedule.asmx/GetScheduleDataX`, {
        params: {
            callback: '',
            aVuzID: 11613,
            aStudyGroupID: `"${studyGroupKey}"`,
            aStartDate: `"${startString}"`,
            aEndDate: `"${endString}"`,
            aStudyTypeID: 'null',
        },
    });
    const data = response.data.d as ScheduleRow[];
    if (!data) throw new Error('Failed to get schedule');

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
        day.classes.push({
            class: row.study_time,
            begin: row.study_time_begin,
            end: row.study_time_end,
            descipline: row.discipline,
            type: row.study_type,
            cabinet: row.cabinet,
            employee: row.employee,
        });
    }

    const days = [...dayByDate.values()].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());

    return days;
}
