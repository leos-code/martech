import service from '../request';

export default {
    get: () => {
        return service({
            url: '/organization/user',
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/organization/user',
            method: 'post',
            data,
        });
    },

};
