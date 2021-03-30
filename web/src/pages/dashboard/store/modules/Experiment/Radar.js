import Vue from 'vue';
import _ from 'lodash';
import { isNotEmptyString } from '@/utils';

const DefaultStage = {
    querying: false,
    options: {},
};

const state = () => ({
    stageMap: {},
});

const getters = {
    getQuerying: (state) => (radarId) => {
        if (isNotEmptyString(radarId) && state.stageMap[radarId]) {
            return state.stageMap[radarId].querying;
        }
        return false;
    },
    queryOptions: (state) => (radarId) => {
        if (isNotEmptyString(radarId) && state.stageMap[radarId]) {
            return state.stageMap[radarId].options;
        }
        return {};
    },
};

const actions = {};

const mutations = {
    setStageMap: (state, { radarId }) => {
        if (isNotEmptyString(radarId) && !state.stageMap[radarId]) {
            Vue.set(state.stageMap, radarId, _.cloneDeep(DefaultStage));
        }
    },
    delStageMap: (state, { radarId }) => {
        if (isNotEmptyString(radarId) && state.stageMap[radarId]) {
            delete state.stageMap[radarId];
        }
    },
    activeQuerying: (state, { radarId, options }) => {
        if (isNotEmptyString(radarId) && state.stageMap[radarId]) {
            state.stageMap[radarId].querying = true;
            state.stageMap[radarId].options = options;
        }
    },
    closeQuerying: (state, { radarId }) => {
        if (isNotEmptyString(radarId) && state.stageMap[radarId]) {
            state.stageMap[radarId].querying = false;
        }
    },
};

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations,
};
