import service from '../request';

export default {
    getList: () => {
        return service({
            url: '/experiment/account/list',
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/experiment/account/edit',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/experiment/account/delete/${id}`,
            method: 'get',
        });
    },
    sync: (data) => {
        return service({
            url: `/experiment/account/sync`,
            method: 'post',
            data,
        });
    },
    getRtaExpList: (data) => {
        return service({
            url: `/experiment/rta_exp/list`,
            method: 'post',
            data,
        });
    },
};
