import { ScheduleDay, getSchedule } from '../api/getSchedule';
import { CurrentKeyboardVersion } from '../const';
import { User } from '../db';
import { askGroup } from './askData';
import { BotContext, bot } from './bot';
import { keyboard } from './commands/commands';


export const escapeMsg = (str: string) => {
    str = str.replace(/\./g, '\\.');
    str = str.replace(/\-/g, '\\-');
    str = str.replace(/\_/g, '\\_');
    str = str.replace(/\*/g, '\\*');
    str = str.replace(/\(/g, '\\(');
    str = str.replace(/\)/g, '\\)');
    str = str.replace(/\!/g, '\\!');
    return str;
}

interface Options {
    forward?: number;
    days?: number;
    startFromMonday?: boolean;
}

const dayToText = async (day: ScheduleDay, sendGroups: boolean) => {
    let message = '';
    message += `${escapeMsg(day.weekday)}, ${escapeMsg(day.date)}\n\n`;
    for (const class1 of day.classes) {
        message += `⚪ *${escapeMsg(class1.class)}*, \\[${escapeMsg(class1.begin)}\\-${escapeMsg(class1.end)}\\]\n`;
        message += `Предмет: ${escapeMsg(class1.descipline)}\n`;
        message += `Вчитель: ${escapeMsg(class1.employee)}\n`;
        message += `Тип: \\[*${escapeMsg(class1.type)}*\\] Кабінет: \\[*${escapeMsg(class1.cabinet)}*\\]\n`;

        if (sendGroups) {
            const groups = await class1.groups;
            if (groups.length > 0) message += `Групи: ${groups.map((value) => escapeMsg(value)).join(', ')}\n`;
        } else {
            message += `Групи: Пошук\\.\\.\\.\n`;
        }

        message += `\n`;
    }

    return message;
}

export const sendSchedule = async (ctx: BotContext, options: Options) => {
    let isKeyboardOutdated = ctx.user.keyboardVersion != CurrentKeyboardVersion;

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

    const schedule = await getSchedule(studyGroupKey, start, end, true);
    if (schedule.length === 0) {
        await ctx.sendMessage('Розклад не знайдено', isKeyboardOutdated ? keyboard : undefined);
        if (isKeyboardOutdated) User.setKeyboardVersion(ctx.user.id, CurrentKeyboardVersion);
        return;
    }

    for (const day of schedule) {
        const messageText = await dayToText(day, false);
        const message = await ctx.replyWithMarkdownV2(messageText, isKeyboardOutdated ? keyboard : undefined);
        if (isKeyboardOutdated) User.setKeyboardVersion(ctx.user.id, CurrentKeyboardVersion);
        isKeyboardOutdated = false;
        dayToText(day, true).then((text) => bot.telegram.editMessageText(ctx.chat!.id, message.message_id, undefined, text, { parse_mode: 'MarkdownV2' }));
    }
}