import { Markup } from 'telegraf';

interface Description {
    command: string;
    description: string;
    startDesc?: string;
    keyboard?: boolean;
}

export const descriptions = [
    { command: 'start', description: 'Старт', keyboard: true } as const,
    { command: 'day', description: 'Розклад на сьогодні', keyboard: true } as const,
    { command: 'next', description: 'Розклад на завтра', keyboard: true, startDesc: 'Розклад на наступні дні(/next@2 - розклад на післязавтра)' } as const,
    { command: 'nextnext', description: 'Розклад на післязавтра', keyboard: true } as const,
    { command: 'week', description: 'Розклад на тиждень', keyboard: true, startDesc: 'Розклад на тиждень(/week@1 - розклад на наступний тиждень)'} as const,
    { command: 'keyboard', description: 'Ввімкнути клавіатуру' } as const,
    { command: 'keyboardoff', description: 'Вимкнути клавіатуру' } as const,
    { command: 'setgroup', description: 'Змінити групу', keyboard: true }as const,
    { command: 'setdata', description: 'Змінити дані', keyboard: true }as const,
] satisfies Description[];

export const keyboard = Markup.keyboard(descriptions
    .filter((value) => value.keyboard)
    .map((value) => value.description), { wrap: (btn, index) => index % 2 === 1 });