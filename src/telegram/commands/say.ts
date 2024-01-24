import { User } from '../../db';
import { command } from './command';

command('say', (ctx) => {
    const match = ctx.message.text.match(/^\/say +([^\s]+) +(.+)$/);
    if (!match) {
        ctx.sendMessage('Не правильні вхідні дані');
        return;
    }

    const username = match[1];
    const userId = parseInt(username);
    const message = match[2];

    if (username === 'all!') {
        for (const user of User.findAll()) ctx.telegram.sendMessage(user.id, message);
        return;
    }
    let user = User.findByUsername(username) || (userId && User.find(userId));
    if (!user) {
        ctx.sendMessage('Не вдалося знайти такого користувача');
        return;
    }

    ctx.telegram.sendMessage(user.id, message);
});