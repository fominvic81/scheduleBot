import { bot } from './telegram/bot';
import { getAllEmployees } from './api/getAllEmployees';

const main = async () => {


    bot.launch();

    // Cache empolyees
    await getAllEmployees();
}

main();
