import axios from 'axios';
import service from '@/utils/request';

const dashboardService = axios.create({
    baseURL: process.env.VUE_APP_API_DASHBOARD || process.env.VUE_APP_API_BASE,
    timeout: 9999,
});

export default (options) => {
    return service(options, dashboardService);
};
