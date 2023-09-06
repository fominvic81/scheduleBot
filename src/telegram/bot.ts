import { Context, Telegraf } from 'telegraf';
import { User } from '@prisma/client';
import { prisma } from '../main';

const token = process.env.TELEGRAM_TOKEN;
if (!token) throw new Error('Telegram bot token is not defined');

export interface BotContext extends Context {
    user: User;
}

export const bot = new Telegraf<BotContext>(token);
bot.telegram.setMyCommands([
    { command: 'start', description: 'Почати' },
    { command: 'setgroup', description: 'Змінити групу' },
    { command: 'setdata', description: 'Змінити дані' },
    { command: 'schedule', description: 'Розклад на поточний тиждень' },
]);

bot.use(async (ctx, next) => {
    if (!ctx.from) return;
    const id = ctx.from.id;
    let user = await prisma.user.findFirst({ where: { id }});

    if (!user) {
        user = await prisma.user.create({
            data: {
                id,
                firstname: ctx.from.first_name,
                lastname: ctx.from.last_name,
                username: ctx.from.username,
            },
        });
    }
    ctx.user = user;
    await next();
})

bot.launch();

import './commands/start';
import './commands/setgroup';
import './commands/setdata';
import './commands/schedule';