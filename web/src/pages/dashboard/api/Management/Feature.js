import service from '../request';

export default {
    get: () => {
        return service({
            url: `/management/feature`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/management/feature',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/management/feature/${id}`,
            method: 'delete',
        });
    },
    relate: (data) => {
        return service({
            url: `/management/feature/relate`,
            method: 'post',
            data,
        });
    },
    sync: () => {
        return service({
            url: `/management/feature/sync`,
            method: 'get',
        });
    },
};
