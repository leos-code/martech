import Vue from 'vue';
import Vuex from 'vuex';

import bindStrategy from './modules/Rta/BindStrategy';
import experimentParameter from './modules/Experiment/Parameter';
import experimentAccount from './modules/Experiment/Account';
import experimentGroup from './modules/Experiment/Group';
import experimentRadar from './modules/Experiment/Radar';
import organizationUser from './modules/Organization/User';
import superAdmin from './modules/Management/SuperAdmin';
import tenant from './modules/Management/Tenant';
import user from './modules/User';

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== 'production';

export default new Vuex.Store({
    modules: {
        bindStrategy,
        user,
        experiment: {
            namespaced: true,
            modules: {
                parameter: experimentParameter,
                account: experimentAccount,
                group: experimentGroup,
                radar: experimentRadar,
            },
        },
        management: {
            namespaced: true,
            modules: {
                superadmin: superAdmin,
                tenant,
            },
        },
        organization: {
            namespaced: true,
            modules: {
                user: organizationUser,
            }
        }
    },
    strict: debug,
});
