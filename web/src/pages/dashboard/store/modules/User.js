import _ from 'lodash';
import Base from '../base';
import Utils from '@/utils';
import { User } from '@dashboard/api';

const state = {
    info: {
        user: {
            id: undefined,
            phone_number: '',
            email: '',
        },
        login_user: {
            open_id: '',
            login_type: '',
            user_id: undefined,
            avatar: '',
            nick_name: '',
        },
        current_tenant: {},
        tenant: [],
    },
    tenantSwitch: {
        active: false,
    },
};

const getters = {
    info(state) {
        return state.info;
    },
    user(state) {
        return state.info.user;
    },
    loginUser(state) {
        return state.info.login_user;
    },
    currentTenant(state) {
        return state.info.current_tenant || {};
    },
    tenant(state) {
        return state.info.tenant;
    },
    tenantSwitchActive(state) {
        return state.tenantSwitch.active;
    },
};

const actions = {
    async initUserInfo(context) {
        if (!Utils.isNum(context?.state?.info?.user?.id)) {
            await context.dispatch('getUserInfo');
        }
    },
    async getUserInfo({ commit }) {
        const { data } = await User.get();
        if (typeof data !== 'undefined') {
            commit('setUserInfo', { info: data });
        }
    },
    activeTenantSwitch({ commit }) {
        commit('activeTenantSwitch');
    },
    closeTenantSwitch({ commit }) {
        commit('closeTenantSwitch');
    },
};

const mutations = {
    setUserInfo: (state, { info }) => {
        state.info = _.assign({}, state.info, info);
    },
    activeTenantSwitch: (state) => {
        state.tenantSwitch.active = true;
    },
    closeTenantSwitch: (state) => {
        state.tenantSwitch.active = false;
    },
};

export default _.merge({}, Base, {
    namespaced: true,
    state,
    getters,
    actions,
    mutations,
});
