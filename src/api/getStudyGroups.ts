import axios from 'axios';
import { KeyValue } from '.';

export const getStudyGroups = async (facultyId: string, educationFormId: string, courseId: string): Promise<KeyValue[]> => {
    
    const response = await axios(`https://vnz.osvita.net/WidgetSchedule.asmx/GetStudyGroups`, {
        params: {
            callback: '',
            aVuzID: 11613,
            aFacultyID: `"${facultyId}"`,
            aEducationForm: `"${educationFormId}"`,
            aCourse: `"${courseId}"`,     
            aGiveStudyTimes: false,
        },
    });
    const data = response.data.d?.studyGroups;
    if (!data) throw new Error('Failed to get study groups');
    
    return data;

}