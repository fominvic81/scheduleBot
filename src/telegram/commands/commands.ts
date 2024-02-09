import { Markup } from 'telegraf';

interface Command {
    command: string;
    description?: string;
    startDesc?: string;
    admin?: boolean;
}

export const commands = [
    { command: 'day', description: 'Розклад на сьогодні' } as const,
    { command: 'next', description: 'Розклад на завтра', startDesc: 'Розклад на наступні дні(/next@2 - розклад на післязавтра)' } as const,
    { command: 'nextnext', description: 'Розклад на післязавтра' } as const,
    { command: 'week', description: 'Розклад на тиждень', startDesc: 'Розклад на тиждень(/week@1 - розклад на наступний тиждень)' } as const,
    { command: 'nextweek', description: 'Розклад на наступний тиждень' } as const,
    { command: 'short', description: 'Стисло(два тижні)', startDesc: 'Розклад на два тижні в компактному форматі' } as const,
    { command: 'start', description: 'Старт' } as const,
    { command: 'setgroup', description: 'Змінити групу' } as const,
    { command: 'setdata', description: 'Змінити дані' } as const,
    { command: 'users', admin: true } as const,
    { command: 'say', admin: true } as const,
] satisfies Command[];

export const publicCommands = commands.filter((command) => !command.admin).map((command) => ({
    command: command.command,
    description: command.description ?? '',
    startDesc: command.startDesc ?? command.description ?? '',
}));

export const keyboard = Markup.keyboard(publicCommands
    .map((value) => value.description), { wrap: (btn, index) => index % 2 === 1 });