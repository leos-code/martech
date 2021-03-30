import service from '../request';

export default {
    getList: (data) => {
        return service({
            url: `/targeting/list?page=${data.page}&page_size=${data.page_size}&filter=${data.filter}`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/targeting/edit',
            method: 'post',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/targeting/delete/${id}`,
            method: 'get',
        });
    },
};
