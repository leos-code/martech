import service from '../request';

export default {
    get: () => {
        return service({
            url: `/management/frontend`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/management/frontend',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/management/frontend/${id}`,
            method: 'delete',
        });
    },
};
