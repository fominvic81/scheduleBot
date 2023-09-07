import { askGroup } from '../askData';
import { command } from './command';

command('setgroup', (ctx) => {
    askGroup(ctx);
});