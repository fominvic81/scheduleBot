import { getAllEmployees } from '../api/getAllEmployees';
import { getEmployeeSchedule } from '../api/getEmployeeSchedule';
import { getSchedule } from '../api/getSchedule';
import { askGroup } from './askData';
import { BotContext } from './bot';


const escapeMsg = (str: string) => {
    str = str.replace(/\./g, '\\.');
    str = str.replace(/\-/g, '\\-');
    str = str.replace(/\_/g, '\\_');
    str = str.replace(/\*/g, '\\*');
    str = str.replace(/\(/g, '\\(');
    str = str.replace(/\)/g, '\\)');
    return str;
}

interface Options {
    forward?: number;
    days?: number;
    startFromMonday?: boolean;
}

export const sendSchedule = async (ctx: BotContext, options: Options) => {
    options.days = options.days ?? 1;
    options.forward = options.forward ?? 0;
    options.startFromMonday = options.startFromMonday ?? false;

    const start = new Date();
    start.setDate(start.getDate() + options.forward);

    if (options.startFromMonday) {
        // if it is sunday get schedule for next week
        if (start.getDay() === 0) start.setDate(start.getDate() + 7);
        
        // Find last monday
        start.setDate(start.getDate() - (start.getDay() + 6) % 7);
    }

    const end = new Date(start);
    end.setDate(end.getDate() + options.days - 1);

    let studyGroup = ctx.user.studyGroup;

    if (!studyGroup) {
        studyGroup = (await askGroup(ctx)).Key;
    }

    const studyGroupKey = studyGroup;

    const schedule = await getSchedule(studyGroupKey, start, end, false);

    for (const day of schedule) {
        let message = '';
        message += `${escapeMsg(day.weekday)}, ${escapeMsg(day.date)}\n\n`;

        for (const class1 of day.classes) {

            message += `⚪ *${escapeMsg(class1.class)}*\n`;
            message += `Час: ${escapeMsg(class1.begin)}\\-${escapeMsg(class1.end)}\n`;
            message += `Предмет: ${escapeMsg(class1.descipline)}\n`;
            message += `Вчитель: ${escapeMsg(class1.employee)}\n`;
            message += `Кабінет: ${escapeMsg(class1.cabinet)}\n`;
            message += `Тип заняття: ${escapeMsg(class1.type)}\n`;

            if (class1.groups.length > 0) message += `Групи: ${class1.groups.map((value) => escapeMsg(value)).join(', ')}\n`;

            message += `\n`;
        }

        await ctx.replyWithMarkdownV2(message);
    }

    if (schedule.length === 0) await ctx.sendMessage('Розклад не знайдено');
}