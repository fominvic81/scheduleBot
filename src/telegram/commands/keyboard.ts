import { Markup } from 'telegraf';
import { descriptions } from './descriptions';
import { command } from './command';

command('keyboard', (ctx) => {
    ctx.sendMessage('Клавіатуру ввімкнено', Markup.keyboard(descriptions
        .filter((value) => value.keyboard)
        .map((value) => value.description), { columns: 2 }));
});
command('keyboardoff', (ctx) => {
    ctx.sendMessage('Клавіатуру вимкнено', Markup.removeKeyboard());
});
