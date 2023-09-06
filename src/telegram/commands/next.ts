import { bot } from '../bot';
import { sendSchedule } from '../sendSchedule';

bot.command('next', (ctx) => {
    let forward = parseInt(ctx.payload);
    if (Number.isNaN(forward)) forward = 0;

    sendSchedule(ctx, { days: 1, forward });
});