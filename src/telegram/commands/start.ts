import { askInfo } from '../askData';
import { keyboard, publicCommands } from './commands';
import { command } from './command';

command('start', async (ctx) => {

    await ctx.reply(
        'Я бот, що дозволяє зручно та швидко слідкувати за розкладом занять в ЛНТУ\n' + 
        'Основні команди: \n' + 
        publicCommands.map(description => `/${description.command} - ${description.startDesc ?? description.description}`).join('\n') + '\n',
        keyboard,
    );

    if (!ctx.user.studyGroup) {
        askInfo(ctx);
    }
});