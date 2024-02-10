import { Markup } from 'telegraf';

interface Command {
    command: string;
    description?: string;
    startDesc?: string;
    admin?: boolean;
}

export const commands = [
    [
        { command: 'day', description: 'Розклад на сьогодні' } as const as Command,
        { command: 'next', description: 'Розклад на завтра', startDesc: 'Розклад на наступні дні(/next@2 - розклад на післязавтра)' } as const as Command,
        { command: 'nextnext', description: 'Розклад на післязавтра' } as const as Command,
    ],
    [
        { command: 'week', description: 'Розклад на тиждень', startDesc: 'Розклад на тиждень(/week@1 - розклад на наступний тиждень)' } as const as Command,
        { command: 'nextweek', description: 'Розклад на наступний тиждень' } as const as Command,
    ],
    [
        { command: 'short', description: 'Стисло(два тижні)', startDesc: 'Розклад на два тижні в компактному форматі' } as const as Command,
        { command: 'subject', description: 'Знайти предмет' } as const as Command,
    ],
    [
        { command: 'setgroup', description: 'Змінити групу' } as const as Command,
        { command: 'start', description: 'Старт' } as const as Command,
        { command: 'setdata', description: 'Змінити дані' } as const as Command,
    ],
    [
        { command: 'users', admin: true } as const as Command,
        { command: 'say', admin: true } as const as Command,
    ],
];

export const publicCommands = commands.map((row) => row.filter((command) => !command.admin).map((command) => ({
    command: command.command,
    description: command.description ?? '',
    startDesc: command.startDesc ?? command.description ?? '',
})));

export const keyboard = Markup.keyboard(publicCommands
    .map((row) => row.map((command) => command.description)));