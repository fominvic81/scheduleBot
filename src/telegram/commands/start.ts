import { bot } from '../bot';
import { askData } from '../askData';

bot.start(async (ctx) => {

    await ctx.reply('Text');

    askData(ctx);
});