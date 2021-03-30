import service from '../request';

export default {
    getSchema: () => {
        return service({
            url: '/schema',
            method: 'get',
        });
    },
};
