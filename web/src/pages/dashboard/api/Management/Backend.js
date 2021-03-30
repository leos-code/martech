import service from '../request';

export default {
    get: () => {
        return service({
            url: `/management/backend`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/management/backend',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/management/backend/${id}`,
            method: 'delete',
        });
    },
};
