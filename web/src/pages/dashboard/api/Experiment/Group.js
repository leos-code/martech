import service from '../request';

export default {
    getList: () => {
        return service({
            url: '/experiment/group/list',
            method: 'get',
        });
    },
    get: (id) => {
        return service({
            url: `/experiment/group/get/${id}`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/experiment/group/edit',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/experiment/group/delete/${id}`,
            method: 'get',
        });
    },
    prompt: (id) => {
        return service({
            url: `/experiment/group/prompt/${id}`,
            method: 'get',
        });
    },
    stop: (id) => {
        return service({
            url: `/experiment/group/stop/${id}`,
            method: 'get',
        });
    },
};
