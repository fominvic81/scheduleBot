import { Markup } from 'telegraf';
import { BotContext, bot } from './bot';
import { button } from './button';
import { KeyValue } from '../api';
import { getFilters } from '../api/getFilters';
import { getStudyGroups } from '../api/getStudyGroups';
import { User } from '../db';

export const ask = (ctx: BotContext, text: string, buttons: KeyValue[], columns = 4) => new Promise<KeyValue>((resolve) => {
    ctx.sendMessage(text,
        Markup.inlineKeyboard(
            buttons.map(({ Key, Value }) => button({ text: Value }, async (buttonCtx) => {
                buttonCtx.deleteMessage();
                resolve({ Key, Value });
            })),
            { columns },
        ),
    );
});

export const askInfo = async (ctx: BotContext) => {
    const filters = await getFilters();

    const message = await ctx.sendMessage('Вкажіть дані:');

    User.reset(ctx.user.id);

    const course = await ask(ctx, 'Курс', filters.courses);
    const educationForm = await ask(ctx, 'Форма Навчання', filters.educForms);
    const faculty = await ask(ctx, 'Факультет', filters.faculties);
    
    User.setInfo(ctx.user.id, faculty.Key, educationForm.Key, course.Key);
    ctx.user.course = course.Key;
    ctx.user.educationForm = educationForm.Key;
    ctx.user.faculty = faculty.Key;

    const studyGroup = await askGroup(ctx, false);

    if (ctx.chat && studyGroup) {
        bot.telegram.editMessageText(ctx.chat.id, message.message_id, undefined,
            `Курс: ${course.Value}\n` +
            `Форма навчання: ${educationForm.Value}\n` +
            `Факультет: ${faculty.Value}\n` +
            `Навчальна група: ${studyGroup.Value}`,
        );
    } else {
        console.error('Study group or chat is undefined')
    }

    return {
        course,
        educationForm,
        faculty,
        studyGroup,
    };
}

export const askGroup = async (ctx: BotContext, replyWithChosen: boolean = true): Promise<KeyValue> => {

    const course = ctx.user.course;
    const educationForm = ctx.user.educationForm;
    const faculty = ctx.user.faculty;

    if (!course || !educationForm || !faculty) {
        const { studyGroup } = await askInfo(ctx);
        return studyGroup;
    }
    
    const studyGroups = await getStudyGroups(faculty, educationForm, course);

    const studyGroup = await ask(ctx, 'Навчальна Група', studyGroups);

    User.setStudyGroup(ctx.user.id, studyGroup.Key);
    ctx.user.studyGroup = studyGroup.Key

    if (replyWithChosen) {
        ctx.sendMessage(`Навчальна група: ${studyGroup.Value}`);
    }

    return studyGroup;
}