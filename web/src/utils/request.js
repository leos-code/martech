import axios from 'axios';
import Utils from './index';

const service = axios.create({
    baseURL: process.env.VUE_APP_API_BASE,
    timeout: 9999,
});

export default (options, externalService) => {
    return new Promise((resolve, reject) => {

        if (process.env.NODE_ENV === 'development'){
            console.log(options);
        }

        externalService = Utils.isFunction(externalService) ? externalService : service;
        externalService(options)
            .then((response) => {
                const { data: dataWrapper } = response || {};
                const { code, msg, data } = dataWrapper || {};
                if (code === 0) {
                    resolve({
                        data,
                    });
                } else {
                    resolve({
                        code,
                        msg,
                    });
                }
            })
            .catch((err) => {
                console.error(err);
                if (err.response) {
                    const { data: dataWrapper } = err.response || {};
                    let { msg } = dataWrapper || {};
                    msg = msg ? msg : err.toString();
                    resolve({
                        code: -1,
                        err,
                        msg,
                    });
                } else {
                    resolve({
                        code: -1,
                        err,
                        msg: '系统内部错误',
                    });
                }
            });
    });
};
