import { sendSchedule } from '../sendSchedule';
import { command } from './command';

command('nextweek', (ctx) => {
    sendSchedule(ctx, { days: 7, forward: 7, startFromMonday: true });
});