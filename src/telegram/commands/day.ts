import { sendSchedule } from '../sendSchedule';
import { command } from './command';

command('day', (ctx) => {
    sendSchedule(ctx, { days: 1 });
});