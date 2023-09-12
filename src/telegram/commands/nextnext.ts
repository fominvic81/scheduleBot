import { sendSchedule } from '../sendSchedule';
import { command } from './command';

command('nextnext', (ctx) => {
    sendSchedule(ctx, { days: 1, forward: 2 });
});