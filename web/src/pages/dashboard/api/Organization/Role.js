import service from '../request';

export default {
    get: () => {
        return service({
            url: '/organization/role',
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/organization/role',
            method: 'post',
            data,
        });
    },
    editUser: (data) => {
        return service({
            url: '/organization/role/user',
            method: 'post',
            data,
        });
    },
    editObject: (data) => {
        return service({
            url: '/organization/role/object',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/organization/role/${id}`,
            method: 'delete',
        });
    },
};
