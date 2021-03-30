import service from '../request';

export default {
    get: () => {
        return service({
            url: `/organization/policy`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/organization/policy',
            method: 'post',
            data,
        });
    },
};
