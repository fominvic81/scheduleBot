import { User } from '../../db';
import { escapeMsg } from '../sendSchedule';
import { command } from './command';

command('users', (ctx) => {
    
    let message = '';
    const users = User.findAll();
    users.sort((a, b) => a.messages - b.messages);
    for (const user of users) {
        let userStr = '';
        userStr += `${user.id} \\| ${escapeMsg(user.username ?? '')}\n`;
        userStr += `${escapeMsg((user.firstname + (user.lastname ?? '')))}\n`;
        userStr += `${user.messages}\n`;
        userStr += '\n';

        if (message.length + userStr.length > 4000) {
            ctx.sendMessage(message, { parse_mode: 'MarkdownV2' });
            message = '';
        }
        message += userStr;
    }

    ctx.sendMessage(message, { parse_mode: 'MarkdownV2' });
});