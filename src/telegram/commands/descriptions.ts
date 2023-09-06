
interface Description {
    command: string;
    description: string;
    startDesc?: string;
}

export const descriptions = [
    { command: 'start', description: 'Старт' },
    { command: 'setgroup', description: 'Встановити групу' },
    { command: 'setdata', description: 'Встановити дані' },
    { command: 'day', description: 'Розклад на сьогодні' },
    { command: 'next', description: 'Розклад на завтра', startDesc: 'Розклад на наступні дні(/next@2 - розклад на післязавтра)' },
    { command: 'week', description: 'Розклад на тиждень', startDesc: 'Розклад на тиждень(/week@1 - розклад на наступний тиждень)'},
] satisfies Description[];