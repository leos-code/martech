import service from '../request';

export default {
    get: (data) => {
        const { type } = data || {};
        return service({
            url: `/organization/object${type ? `?type=${type}` : ''}`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/organization/object',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/organization/object/${id}`,
            method: 'delete',
        });
    },
};
