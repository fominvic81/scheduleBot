import { Context, Telegraf } from 'telegraf';
import { User } from '@prisma/client';
import { prisma } from '../main';
import { descriptions } from './commands/descriptions';

const token = process.env.TELEGRAM_TOKEN;
if (!token) throw new Error('Telegram bot token is not defined');

export interface BotContext extends Context {
    user: User;
}

export const bot = new Telegraf<BotContext>(token);

bot.telegram.setMyCommands(descriptions);
bot.telegram.setMyDescription('Розклад для студентів лнту');

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
});

bot.catch((err, ctx) => {
    console.error(err, ctx);
});

import './commands/start';
import './commands/setgroup';
import './commands/setdata';
import './commands/day';
import './commands/next';
import './commands/week';
import './commands/nextnext';
import './commands/keyboard';
