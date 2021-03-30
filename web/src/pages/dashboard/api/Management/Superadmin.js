import service from '../request';

export default {
    get: (id) => {
        return service({
            url: `/management/superadmin`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/management/superadmin',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/management/superadmin/${id}`,
            method: 'delete',
        });
    },
};
