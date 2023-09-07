import { askData } from '../askData';
import { command } from './command';

command('setdata', (ctx) => {
    askData(ctx);
});