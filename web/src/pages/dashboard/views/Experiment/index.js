export const Experiment = () => import('./index.vue');

export const ExperimentAccount = () => import('./Account/index.vue');
export const ExperimentAccountList = () => import('./Account/list');
export const ExperimentAccountEdit = () => import('./Account/edit');

export const ExperimentParameter = () => import('./Parameter/index.vue');
export const ExperimentParameterList = () => import('./Parameter/list');
export const ExperimentParameterEdit = () => import('./Parameter/edit');

export const ExperimentGroup = () => import('./Group/index.vue');
export const ExperimentGroupList = () => import('./Group/list.vue');
export const ExperimentGroupEdit = () => import('./Group/edit.vue');
export const ExperimentGroupDetail = () => import('./Group/Detail/index.vue');
export const ExperimentGroupDetailPage = () => import('./Group/Detail/page.vue');
export const ExperimentGroupDetailDraft = () => import('./Group/Detail/Draft');
export const ExperimentGroupDetailCurrent = () => import('./Group/Detail/Current');
export const ExperimentGroupDetailRecord = () => import('./Group/Detail/Record');
export const ExperimentGroupDetailRadar = () => import('./Group/Detail/Radar');
