import { bot } from '../bot';
import { sendSchedule } from '../sendSchedule';

bot.command('week', (ctx) => {
    let forward = parseInt(ctx.payload);
    if (Number.isNaN(forward)) forward = 0;
    
    forward *= 7;

    sendSchedule(ctx, { days: 7, forward, startFromMonday: forward > 0 });
});