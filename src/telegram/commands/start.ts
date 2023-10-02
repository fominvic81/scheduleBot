import { askInfo } from '../askData';
import { descriptions, keyboard } from './descriptions';
import { command } from './command';

command('start', async (ctx) => {

    await ctx.reply(
        'Я бот, що дозволяє зручно та швидко слідкувати за розкладом занять в ЛНТУ\n' + 
        'Основні команди: \n' + 
        descriptions.map(description => `/${description.command} - ${description.startDesc ?? description.description}`).join('\n') + '\n',
        keyboard,
    );

    if (!ctx.user.studyGroup) {
        askInfo(ctx);
    }
});