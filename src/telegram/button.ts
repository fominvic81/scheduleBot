import { Markup } from 'telegraf';
import { bot, BotContext } from './bot';

type Data = Parameters<typeof Markup.inlineKeyboard>[0][number];
type DataWithoutCallback = Omit<Data, 'callback_data'>;
type Callback = (context: BotContext) => any;

const callbacks: Map<number, Callback> = new Map();
let lastId = 0;

export const button = (data: DataWithoutCallback, callback: Callback): Data => {
    const id = lastId++;
    callbacks.set(id, callback);

    return {
        ...data,
        callback_data: `btn:${id}`,
    }
}

bot.action(/btn:(.*)/, (ctx) => {

    const id = parseInt(ctx.match[1]);
    const callback = callbacks.get(id);

    if (!callback) {
        console.error('Unknow button was pressed');
        return;
    }

    callback(ctx);
});