import { Markup } from 'telegraf';
import { KeyValue, getFiltersData, getStudyGroupByFilters } from '../scheduleApi';
import { BotContext, bot } from './bot';
import { prisma } from '../main';
import { button } from './button';

const ask = (ctx: BotContext, text: string, buttons: KeyValue[]) => new Promise<KeyValue>((resolve) => {
    ctx.sendMessage(text,
        Markup.inlineKeyboard(
            buttons.map(({ Key, Value }) => button({ text: Value }, async (buttonCtx) => {

                buttonCtx.deleteMessage();
                resolve({ Key, Value });
            })),
            { columns: 2 },
        ),
    );
});

export const askData = async (ctx: BotContext) => {
    const filters = await getFiltersData();

    const message = await ctx.sendMessage('Вкажіть дані:');

    await prisma.user.update({
        where: { id: ctx.user.id },
        data: {
            course: null,
            educationForm: null,
            faculty: null,
            studyGroup: null,
        },
    });

    const course = await ask(ctx, 'Курс', filters.courses.map(value => ({ ...value })));
    const educationForm = await ask(ctx, 'Форма Навчання', filters.educForms.map(value => ({ ...value })));
    const faculty = await ask(ctx, 'Факультет', filters.faculties.map(value => ({ ...value })));
    
    ctx.user = await prisma.user.update({
        where: { id: ctx.user.id },
        data: {
            course: course.Key,
            educationForm: educationForm.Key,
            faculty: faculty.Key,
        },
    });

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
        const { studyGroup } = await askData(ctx);
        return studyGroup;
    }

    const studyGroups = await getStudyGroupByFilters(faculty, educationForm, course);

    const studyGroup = await ask(ctx, 'Навчальна Група', studyGroups.map(value => ({ ...value })));

    ctx.user = await prisma.user.update({
        where: { id: ctx.user.id },
        data: {
            studyGroup: studyGroup as {},
        },
    });

    if (replyWithChosen) {
        ctx.sendMessage(`Навчальна група: ${studyGroup.Value}`);
    }

    return studyGroup;
}