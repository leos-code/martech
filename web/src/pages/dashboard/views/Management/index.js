export const Management = () => import('./index.vue');

export const SuperAdmin = () => import('./SuperAdmin/index');
export const SuperAdminList = () => import('./SuperAdmin/list');
export const SuperAdminEdit = () => import('./SuperAdmin/edit');

export const Tenant = () => import('./Tenant/index');
export const TenantList = () => import('./Tenant/list');
export const TenantEdit = () => import('./Tenant/edit');
export const TenantSwitch = () => import('./Tenant/switch');

export const Frontend = () => import('./Frontend/index');
export const FrontendList = () => import('./Frontend/list');
export const FrontendEdit = () => import('./Frontend/edit');

export const Backend = () => import('./Backend/index');
export const BackendList = () => import('./Backend/list');
export const BackendEdit = () => import('./Backend/edit');

export const Feature = () => import('./Feature/index');
export const FeatureList = () => import('./Feature/list');
export const FeatureEdit = () => import('./Feature/edit');
