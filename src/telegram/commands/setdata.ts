import { askInfo } from '../askData';
import { command } from './command';

command('setdata', (ctx) => {
    askInfo(ctx);
});