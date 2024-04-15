import axios from 'axios';
import { KeyValue } from '.';

export interface FiltersData {
    __type: string;
    faculties: KeyValue[];
    educForms: KeyValue[];
    courses: KeyValue[];
}

export const getFilters = async (): Promise<FiltersData> => {
    
    const response = await axios(`https://vnz.osvita.net/WidgetSchedule.asmx/GetStudentScheduleFiltersData`, {
        params: {
            callback: '',
            aVuzID: 11613,
        },
    });
    
    const data = response.data.d;
    if (!data) throw new Error('Failed to get filters');

    return data;
}