import service from './request';

export default {
    oauth: () => {
        return service({
            url: '/login/oauth',
            method: 'get',
        });
    },
    logout: () => {
        return service({
            url: '/logout',
            method: 'get',
        });
    },
};
