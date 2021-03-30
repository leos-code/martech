import service from './request';

export default {
    get: () => {
        return service({
            url: '/user/info',
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/user/info',
            method: 'post',
            data,
        });
    },
    tenant: (data) => {
        return service({
            url: '/user/tenant',
            method: 'post',
            data,
        });
    },
    authority: () => {
        return service({
            url: '/user/authority',
            method: 'get',
        });
    },
    search: (data) => {
        return service({
            url: '/user/search',
            method: 'post',
            data,
        });
    },
};
