import { getSchedule } from '../../api/getSchedule';
import { CurrentKeyboardVersion } from '../../const';
import { User } from '../../db';
import { escapeMsg, sendSchedule } from '../sendSchedule';
import { command } from './command';
import { keyboard } from './commands';

command('short', async (ctx) => {
    let isKeyboardOutdated = ctx.user.keyboardVersion != CurrentKeyboardVersion;

    const start = new Date();
    const end = new Date();
    end.setDate(end.getDate() + 14);

    const schedule = await getSchedule(ctx.user.studyGroup, start, end, false);
    
    for (const day of schedule) {
        let message = '';
        message += `\n${escapeMsg(day.weekday)}, ${escapeMsg(day.date)}\n`;
        for (const class1 of day.classes) {
            message += `${escapeMsg(class1.class).match(/\d+/)![0]}: ${escapeMsg(class1.descipline)}, ${escapeMsg(class1.type)}\n`;
        }
        await ctx.replyWithMarkdownV2(message, isKeyboardOutdated ? keyboard : undefined);
        if (isKeyboardOutdated) User.setKeyboardVersion(ctx.user.id, CurrentKeyboardVersion);
        isKeyboardOutdated = false;
    }
});