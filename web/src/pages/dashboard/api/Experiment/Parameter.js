import service from '../request';

export default {
    getList: () => {
        return service({
            url: '/experiment/parameter/list',
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/experiment/parameter/edit',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/experiment/parameter/delete/${id}`,
            method: 'get',
        });
    },
};
