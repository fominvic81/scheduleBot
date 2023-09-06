import { getSchedule } from '../scheduleApi';
import { askGroup } from './askData';
import { BotContext } from './bot';


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

    const schedule = await getSchedule(studyGroupKey, start, end);

    for (const day of schedule) {
        let message = '';
        message += `${day.weekday}, ${day.date}\n`;

        for (const class1 of day.classes) {
            message += `\n`;
            message += `*${class1.class}*\n`;
            message += `Час: ${class1.begin}-${class1.end}\n`;
            message += `Предмет: ${class1.descipline}\n`;
            message += `Вчитель: ${class1.employee}\n`;
            message += `Тип заняття: ${class1.type}\n`;
        }

        message = message.replace(/\./g, '\\.');
        message = message.replace(/\-/g, '\\-');

        await ctx.replyWithMarkdownV2(message);
    }

    if (schedule.length === 0) await ctx.sendMessage('Розклад не знайдено');
}