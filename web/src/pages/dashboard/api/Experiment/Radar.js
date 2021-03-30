import service from '../request';

export default {
    get: (data) => {
        return service({
            url: '/experiment/report/get',
            method: 'POST',
            data: data,
        });
    },
};
