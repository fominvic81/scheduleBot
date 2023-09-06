import { bot } from '../bot';
import { askData } from '../askData';
import { descriptions } from './descriptions';

bot.start(async (ctx) => {

    await ctx.reply(
        'Я бот, що дозволяє зручно та швидко слідкувати за розкладом занять в ЛНТУ\n' + 
        'Основні команди: \n' + 
        descriptions.map(description => `/${description.command} - ${description.startDesc ?? description.description}`).join('\n') + '\n',
    );

    if (!ctx.user.studyGroup) {
        askData(ctx);
    }
});