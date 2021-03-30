import { superAdminTenantId, localUserId, localRoleId, localTenantId } from './meta';

export const casbinMenuAdmin = {
    code: 0,
    msg: 'get success',
    data: JSON.stringify({
        m:
            '\n                [request_definition]\n                r = sub, dom, obj, act\n                [policy_definition]\n                p = sub, dom, obj, act\n                [role_definition]\n                g = _, _, _\n                g2 = _, _, _\n                [policy_effect]\n                e = some(where (p.eft == allow))\n                [matchers]\n                m = g(r.sub, p.sub, r.dom) && g2(r.obj, p.obj, r.dom) && r.dom == p.dom && r.act == p.act\n            ',
        p: [
            // SuperTenant
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'AdReport', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'Material', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'Advertiser', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'ExperimentGroup', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'ExperimentParameter', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'ExperimentAccount', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'SuperAdmin', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'Tenant', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'Frontend', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'Backend', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'OrganizationInfo', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'User', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'OrganizationObject', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'Role', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${superAdminTenantId}`, 'OrganizationPolicy', 'read'],
            ['p', 'designer', `tenant_${superAdminTenantId}`, 'Material', 'read'],
            ['p', 'designer', `tenant_${superAdminTenantId}`, 'AdReport', 'read'],
            ['p', 'agent', `tenant_${superAdminTenantId}`, 'Advertiser', 'read'],
            ['p', 'agent', `tenant_${superAdminTenantId}`, 'AdReport', 'read'],
            ['p', 'analyst', `tenant_${superAdminTenantId}`, 'ExperimentGroup', 'read'],
            ['p', 'analyst', `tenant_${superAdminTenantId}`, 'ExperimentParameter', 'read'],
            ['p', 'analyst', `tenant_${superAdminTenantId}`, 'ExperimentAccount', 'read'],
            ['p', 'analyst', `tenant_${superAdminTenantId}`, 'AdReport', 'read'],
            //
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'AdReport', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'Material', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'Advertiser', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'ExperimentGroup', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'ExperimentParameter', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'ExperimentAccount', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'SuperAdmin', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'Tenant', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'Frontend', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'Backend', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'OrganizationInfo', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'User', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'OrganizationObject', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'Role', 'read'],
            ['p', `role_${localRoleId}`, `tenant_${localTenantId}`, 'OrganizationPolicy', 'read'],
        ],
        g: [
            ['g', 'user_2', `role_${localRoleId}`, `tenant_${superAdminTenantId}`],
            ['g', 'crystalye', `role_${localRoleId}`, `tenant_${superAdminTenantId}`],
            ['g', 'slyao', 'designer', `tenant_${superAdminTenantId}`],
            ['g', `user_${localUserId}`, `role_${localRoleId}`, `tenant_${superAdminTenantId}`],
            ['g', `user_${localUserId}`, `role_${localRoleId}`, `tenant_${localTenantId}`],
            ['g', 'user_2', `role_${localRoleId}`, `tenant_${localTenantId}`],
            ['g', 'lynlinpeng', 'agent', `tenant_${superAdminTenantId}`],
            ['g', 'v_frostyang', 'agent', `tenant_${superAdminTenantId}`],
            ['g', 'freyaxtzhou', 'analyst', `tenant_${superAdminTenantId}`],
            ['g', 'zelfaliu', 'analyst', `tenant_${superAdminTenantId}`],
        ],
    }),
};

export const casbinEmptyData = {
    code: 0,
    msg: 'get success',
};
