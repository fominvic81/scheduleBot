import { Context, Telegraf } from 'telegraf';
import { publicCommands } from './commands/commands';
import { User } from '../db';
import { UserI } from '../db/types';
import { Message, Update } from 'telegraf/types';

const token = Bun.env.TELEGRAM_TOKEN;
if (!token) throw new Error('Telegram bot token is not defined');

export interface BotContext extends Context<{
    message: Update.New & Update.NonChannel & Message.TextMessage,
    update_id: number,
}> {
    user: UserI;
}

export const bot = new Telegraf<BotContext>(token);

bot.telegram.setMyCommands(publicCommands.flat());
bot.telegram.setMyDescription('Розклад для студентів лнту');

bot.use(async (ctx, next) => {

    if (!ctx.from) return;
    const id = ctx.from.id;

    const user = User.findOrCreate(id, ctx.from.first_name, ctx.from.last_name, ctx.from.username);
    ctx.user = user;

    if (ctx.updateType == 'message') {
        User.incrementMessages(user.id);
    }

    await next();
});

bot.catch((err, ctx) => {
    console.error(err, ctx);
});