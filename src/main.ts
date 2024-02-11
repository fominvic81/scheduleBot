import { bot } from './telegram';
import { getAllEmployees } from './api/getAllEmployees';

const main = async () => {
    // Cache empolyees
    getAllEmployees();

    while (true) await bot.launch().catch(console.error);
}

main();
