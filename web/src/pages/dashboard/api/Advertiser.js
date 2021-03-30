import service from './request';

export default {
    get: () => {
        return service({
            url: '/advertiser',
            method: 'get',
        });
    },
    edit: (data) => {
        return service({
            url: '/advertiser',
            method: 'patch',
            data,
        });
    },
    delete: (id) => {
        return service({
            url: `/advertiser/${id}`,
            method: 'delete',
        });
    },
    authorize: () => {
        return service({
            url: `/advertiser/authorize`,
            method: 'get',
        });
    },
};
