import service from './request';

export default {
    get: () => {
        return service({
            url: '/register',
            method: 'get',
        });
    },
    create: (data) => {
        return service({
            url: '/register',
            method: 'post',
            data,
        });
    },
};
