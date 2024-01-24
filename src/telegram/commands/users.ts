import { User } from '../../db';
import { escapeMsg } from '../sendSchedule';
import { command } from './command';

command('users', (ctx) => {
    
    let message = '';
    const users = User.findAll();
    users.sort((a, b) => a.messages - b.messages);
    for (const user of users) {
        
        message += `${user.id} \\| ${escapeMsg(user.username ?? '')}\n`;
        message += `${escapeMsg((user.firstname + (user.lastname ?? '')))}\n`;
        message += `${user.messages}\n`;
        message += '\n';
    }
    message += '';

    ctx.sendMessage(message, { parse_mode: 'MarkdownV2' });

});