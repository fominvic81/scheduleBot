import { KeyValue, getSchedule } from '../../scheduleApi';
import { askData, askGroup } from '../askData';
import { bot } from '../bot';

bot.command('schedule', (ctx) => {
    (async () => {

        let studyGroup = ctx.user.studyGroup;
    
        if (!studyGroup) {
            studyGroup = (await askGroup(ctx)).Key;
        }
    
        const studyGroupKey = studyGroup;
    
        const schedule = await getSchedule(studyGroupKey);
    
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
    })();

});