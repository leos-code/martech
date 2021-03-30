import service from '../request';

export default {
    getList: () => {
        return service({
            url: '/bind_strategy/list',
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/bind_strategy/edit',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/bind_strategy/delete/${id}`,
            method: 'get',
        });
    },
};
