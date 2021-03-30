import service from '../request';

export default {
    get: () => {
        return service({
            url: `/management/tenant`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/management/tenant',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/management/tenant/${id}`,
            method: 'delete',
        });
    },
};
