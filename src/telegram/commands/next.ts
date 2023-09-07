import { sendSchedule } from '../sendSchedule';
import { command } from './command';

command('next', (ctx) => {
    let forward = parseInt(ctx.payload);
    if (Number.isNaN(forward)) forward = 1;

    sendSchedule(ctx, { days: 1, forward });
});