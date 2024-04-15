import axios from 'axios';
import { KeyValue } from '.';


export const getFaculies = async (): Promise<KeyValue[]> => {

    const response = await axios(`https://vnz.osvita.net/WidgetSchedule.asmx/GetEmployeeFaculties`, {
        params: {
            callback: '',
            aVuzID: 11613,
        },
    });
    const data = response.data.d;
    if (!data) throw new Error('Failed to get filters');

    return data;
}