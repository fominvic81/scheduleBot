import { PrismaClient } from '@prisma/client';
import { bot } from './telegram/bot';
import { getAllEmployees } from './api/getAllEmployees';

export const prisma = new PrismaClient();

const main = async () => {
    await prisma.$connect();

    bot.launch();

    // Cache empolyees
    await getAllEmployees();
}

main()
    .then(async () => {
        await prisma.$disconnect();
    })
    .catch(async (e) => {
        console.error(e);
        await prisma.$disconnect();
        process.exit(1);
    })
