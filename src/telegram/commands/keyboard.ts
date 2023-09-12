import { Markup } from 'telegraf';
import { keyboard } from './descriptions';
import { command } from './command';

command('keyboard', (ctx) => {
    ctx.sendMessage('Клавіатуру ввімкнено', keyboard);
});
command('keyboardoff', (ctx) => {
    ctx.sendMessage('Клавіатуру вимкнено', Markup.removeKeyboard());
});
