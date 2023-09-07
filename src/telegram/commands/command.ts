import { BotContext, bot } from '../bot'
import { descriptions } from './descriptions'

type WithPayload<T> = T & { payload: string };
type Callback = (ctx: WithPayload<BotContext>) => any;

export const command = (command: typeof descriptions[number]['command'], callback: Callback) => {

    bot.command(command, async (ctx) => await callback(ctx));
    const info = descriptions.find((value) => value.keyboard && value.command == command);
    if (info?.keyboard) {
        bot.hears(info.description, async (ctx) => {
            (ctx as WithPayload<typeof ctx>).payload = '';
            await callback(ctx as WithPayload<typeof ctx>);
        });
    }

    return bot;
}