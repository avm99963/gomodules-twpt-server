import configDev from './config.dev.json5';
import configProd from './config.json5';

const config = PRODUCTION ? configProd : configDev;
export default config;
