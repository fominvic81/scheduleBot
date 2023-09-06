import { bot } from '../bot';
import { sendSchedule } from '../sendSchedule';

bot.command('day', (ctx) => {
    sendSchedule(ctx, { days: 1 });
});