import service from '../request';

export default {
    get: () => {
        return service({
            url: '/organization/info',
            method: 'get',
        });
    },
};
