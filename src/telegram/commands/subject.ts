import { getSchedule } from '../../api/getSchedule';
import { CurrentKeyboardVersion } from '../../const';
import { User } from '../../db';
import { ask, askInfo } from '../askData';
import { dayToText, escapeMsg } from '../sendSchedule';
import { command } from './command';
import { keyboard } from './commands';

command('subject', (ctx) => {
    (async () => {
        let isKeyboardOutdated = ctx.user.keyboardVersion != CurrentKeyboardVersion;
    
        const start = new Date();
        const end = new Date();
        end.setDate(end.getDate() + 14);
    
        const schedule = await getSchedule(ctx.user.studyGroup, start, end, false);
        const disciplines = [...new Set(schedule.map((day) => day.classes).flat().map((class1) => class1.descipline))].sort();
        const { Value: discipline } = await ask(ctx, 'Виберіть предмет', disciplines.map((value) => ({ Key: value, Value: value })), 2);
        
        for (const day of schedule) {
            day.classes = day.classes.filter((class1) => class1.descipline == discipline);
            if (day.classes.length == 0) continue;
            const messageText = await dayToText(day, false, false);
            await ctx.replyWithMarkdownV2(messageText, isKeyboardOutdated ? keyboard : undefined);
            if (isKeyboardOutdated) User.setKeyboardVersion(ctx.user.id, CurrentKeyboardVersion);
            isKeyboardOutdated = false;
        }
    })();
});