import { sendSchedule } from '../sendSchedule';
import { command } from './command';

command('week', (ctx) => {
    let forward = parseInt(ctx.payload);
    if (Number.isNaN(forward)) forward = 0;
    
    forward *= 7;

    sendSchedule(ctx, { days: 7 - (forward == 0 ? new Date().getDay() : 0), forward, startFromMonday: forward > 0 });
});