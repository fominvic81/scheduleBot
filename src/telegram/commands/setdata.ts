import { bot } from '../bot';
import { askData } from '../askData';

bot.command('setdata', (ctx) => {
    askData(ctx);
});