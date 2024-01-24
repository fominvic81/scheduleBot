import { BotContext, bot } from '../bot'
import { commands } from './commands'

type WithPayload<T> = T & { payload: string };
type Callback = (ctx: WithPayload<BotContext>) => any;

export const command = (command: typeof commands[number]['command'], callback: Callback) => {

    const info = commands.find((value) => value.command == command);
    if (!info) return console.error(`Command ${command} not found`);

    bot.command(command, async (ctx) => {
        if (info.admin && !ctx.user.isAdmin) return;
        await callback(ctx);
    });
    if (info?.description) {
        bot.hears(info.description, async (ctx) => {
            if (info.admin && !ctx.user.isAdmin) return;
            (ctx as WithPayload<typeof ctx>).payload = '';
            await callback(ctx as WithPayload<typeof ctx>);
        });
    }
}