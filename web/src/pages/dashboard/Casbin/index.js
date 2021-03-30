import Store from '@dashboard/store';
import Utils from '@/utils';
const casbinjs = require('./casbin');

const STATUS = {
    Idle: 0,
    Initing: 1,
    Inited: 2,
};

let instance = null;
let status = STATUS.Idle;
const getResolveQueue = [];

const callbackResolveQueue = (instance) => {
    if (Utils.isNonEmptyArray(getResolveQueue)) {
        getResolveQueue.forEach((resolve) => {
            if (Utils.isFunction(resolve)) {
                resolve(instance);
            }
        });
        getResolveQueue.length = 0;
    }
};

const create = async () => {
    if (status !== STATUS.Idle) {
        return;
    }
    status = STATUS.Initing;
    try {
        const casbin = new casbinjs.Authorizer('auto', {
            endpoint: `${process.env.VUE_APP_API_DASHBOARD}/user/authority`,
        });
        await Store.dispatch('user/initUserInfo');
        const { id: userId } = Store.getters['user/user'];
        if (Utils.isPositiveInteger(userId)) {
            await casbin.setUser(`user_${userId}`);
            instance = casbin;
            status = STATUS.Inited;
        } else {
            console.error(`Casbin create error. userId is invalid`);
            status = STATUS.Idle;
        }
        callbackResolveQueue(instance);
    } catch (error) {
        console.error(`Casbin create error.`);
        console.error(error);
        status = STATUS.Idle;
    }
};

const get = () => {
    return new Promise((resolve) => {
        if (instance) {
            resolve(instance);
        } else {
            getResolveQueue.push(resolve);
            create();
        }
    });
};

const can = async ({ object, action, domain }) => {
    const casbin = await get();
    if (Utils.isNull(casbin)) {
        console.error('Casbin can error: Instance is null');
        return false;
    }
    if (!Utils.isNotEmptyString(object) || !Utils.isNotEmptyString(action)) {
        console.error('Casbin can error: Lack of object or action');
        return false;
    }
    if (!Utils.isNotEmptyString(domain)) {
        const { id: currentTenantId } = Store.getters['user/currentTenant'];
        if (Utils.isPositiveInteger(currentTenantId)) {
            domain = `tenant_${currentTenantId}`;
        } else {
            console.error('Casbin can error: Domain & tenant is empty');
            return false;
        }
    } else {
        if (domain.indexOf('tenant_') < 0) {
            domain = `tenant_${domain}`;
        }
    }
    // console.log(`casbin can: user=${instance.user}, domain=${domain}, object=${object}, action=${action}. result: ${await casbin.can(action, object, domain)}`);
    return await casbin.can(action, object, domain);
};

const clear = async () => {
    const casbin = await get();
    casbin.removeCache();
    instance = null;
};

export default {
    can,
    clear,
};
