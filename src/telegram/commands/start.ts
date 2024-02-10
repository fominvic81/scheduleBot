import { askInfo } from '../askData';
import { keyboard, publicCommands } from './commands';
import { command } from './command';
import { CurrentKeyboardVersion } from '../../const';
import { User } from '../../db';

command('start', async (ctx) => {

    await ctx.reply(
        'Я бот, що дозволяє зручно та швидко слідкувати за розкладом занять в ЛНТУ\n' + 
        'Основні команди: \n' + 
        publicCommands.flat().map(description => `/${description.command} - ${description.startDesc ?? description.description}`).join('\n') + '\n',
        keyboard,
    );
    User.setKeyboardVersion(ctx.user.id, CurrentKeyboardVersion);

    if (!ctx.user.studyGroup) {
        askInfo(ctx);
    }
});