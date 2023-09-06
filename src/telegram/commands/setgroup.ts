import { bot } from '../bot';
import { askGroup } from '../askData';

bot.command('setgroup', (ctx) => {
    askGroup(ctx);
});