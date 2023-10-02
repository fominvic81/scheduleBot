import { bot } from './telegram';
import { getAllEmployees } from './api/getAllEmployees';

const main = async () => {
    bot.launch();

    // Cache empolyees
    await getAllEmployees();
}

main();
