import service from '../request';

export default {
    get: (data) => {
        return service({
            url: `/material/list?page=${data.page}&page_size=${data.page_size}&filter=${data.filter}`,
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: `/material/edit`,
            method: 'post',
            data,
        });
    },
    batchAudit: (data) => {
        return service({
            url: '/material/audit/submit',
            method: 'post',
            data,
        });
    },
    batchDelete: (data) => {
        return service({
            url: '/material/delete_many',
            method: 'post',
            data,
        });
    },
};
