import Vue from 'vue';
import VueRouter from 'vue-router';
import _ from 'lodash';

import Casbin from '@dashboard/Casbin';
import { Home } from '@dashboard/views/Home';
import { Rta, BindStrategy, Targeting, Crowd } from '@dashboard/views/Rta';
import {
    Experiment,
    ExperimentAccount,
    ExperimentParameter,
    ExperimentGroup,
    ExperimentGroupEdit,
    ExperimentGroupDetail,
} from '@dashboard/views/Experiment';
import { Report, AdReport } from '@dashboard/views/Report';
import { Property, Material } from '@dashboard/views/Property';
import { Advertiser, AdvertiserProxy } from '@dashboard/views/Advertiser';
import { Management, SuperAdmin, Tenant, Frontend, Backend, Feature } from '@dashboard/views/Management';
import {
    Organization,
    OrganizationInfo,
    OrganizationObject,
    Role,
    OrganizationPolicy,
    User,
} from '@dashboard/views/Organization';

Vue.use(VueRouter);

const routes = [
    /**
     * 路由和菜单的配置, 合并在同个配置中
     * {
     *      path: '',
     *      name: '',
     *      meta: {
     *          show: true,                 // 菜单是否展示, 默认为false
     *          rights: false,              // 路由和菜单是否做校验, 默认为true
     *          label: '',                  // 菜单文案
     *          icon: '',                   // 菜单icon
     *      },
     *      component: () => import(),  // 渲染组件
     * }
     */

    {
        path: '/home',
        name: 'Home',
        meta: {
            show: true,
            rights: false,
            label: '首页',
            icon: 'el-icon-s-home',
        },
        component: Home,
    },
    // Rta管理
    {
        path: '/rta',
        name: 'Rta',
        meta: {
            show: true,
            label: 'Rta策略管理',
            icon: 'el-icon-s-platform',
        },
        component: Rta,
        children: [
            // 定向
            {
                path: 'targeting',
                name: 'RtaTargeting',
                meta: {
                    show: true,
                    label: 'Rta策略管理',
                    icon: 'el-icon-menu',
                },
                component: Targeting,
            },
            // 策略
            // {
            //     path: '/bind_strategy',
            //     name: 'BindStrategy',
            //     meta: {
            //         show: true,
            //         label: 'BindStrategy',
            //         icon: 'el-icon-menu',
            //     },
            //     component: BindStrategy,
            // },
            {
                path: 'crowd',
                name: 'RtaCrowd',
                meta: {
                    show: true,
                    label: '人群包管理',
                    icon: 'el-icon-menu',
                },
                component: Crowd,
            },
        ],
    },
    // 实验系统
    {
        path: '/experiment',
        name: 'Experiment',
        meta: {
            show: true,
            label: '实验系统',
            icon: 'el-icon-s-platform',
        },
        component: Experiment,
        children: [
            // 实验组
            {
                path: 'group',
                name: 'ExperimentGroup',
                meta: {
                    show: true,
                    label: '实验列表',
                    icon: 'el-icon-menu',
                },
                component: ExperimentGroup,
            },
            {
                path: 'group/edit/:id',
                name: 'ExperimentGroupEdit',
                component: ExperimentGroupEdit,
            },
            {
                path: 'group/detail/:id',
                name: 'ExperimentGroupDetail',
                component: ExperimentGroupDetail,
            },
            // 实验参数
            {
                path: 'parameter',
                name: 'ExperimentParameter',
                meta: {
                    show: true,
                    label: '参数配置',
                    icon: 'el-icon-setting',
                },
                component: ExperimentParameter,
            },
            // 实验RTA账户
            {
                path: 'account',
                name: 'ExperimentAccount',
                meta: {
                    show: true,
                    label: 'RTA账户',
                    icon: 'el-icon-s-custom',
                },
                component: ExperimentAccount,
            },
        ],
    },
    // 数据报表
    {
        path: '/report',
        name: 'Report',
        meta: {
            show: true,
            label: '数据报表',
            icon: 'el-icon-s-data',
        },
        component: Report,
        children: [
            {
                path: 'adreport',
                name: 'AdReport',
                meta: {
                    show: true,
                    label: '广告报表',
                    icon: 'el-icon-s-marketing',
                },
                component: AdReport,
            },
        ],
    },
    // 资产管理
    {
        path: '/property',
        name: 'Property',
        meta: {
            show: true,
            label: '资产管理',
            icon: 'el-icon-folder-opened',
        },
        component: Property,
        children: [
            {
                path: 'material',
                name: 'Material',
                meta: {
                    show: true,
                    label: '素材库',
                    icon: 'el-icon-picture',
                },
                component: Material,
            },
        ],
    },
    // 账号管家
    {
        path: '/advertiser',
        name: 'Advertiser',
        meta: {
            show: true,
            label: '账号管家',
            icon: 'el-icon-s-goods',
        },
        component: Advertiser,
    },
    {
        path: '/advertiser/proxy',
        name: 'AdvertiserProxy',
        meta: {
            rights: false,
        },
        component: AdvertiserProxy,
    },
    // 超级管理员
    {
        path: '/management',
        name: 'Management',
        meta: {
            show: true,
            label: '超级管理员',
            icon: 'el-icon-s-tools',
        },
        component: Management,
        children: [
            {
                path: 'superadmin',
                name: 'SuperAdmin',
                meta: {
                    show: true,
                    label: '管理员列表',
                    icon: 'el-icon-user-solid',
                },
                component: SuperAdmin,
            },
            {
                path: 'tenant',
                name: 'Tenant',
                meta: {
                    show: true,
                    label: '租户管理',
                    icon: 'el-icon-document-copy',
                },
                component: Tenant,
            },
            {
                path: 'feature',
                name: 'Feature',
                meta: {
                    show: true,
                    label: '功能管理',
                    icon: 'el-icon-lock',
                },
                component: Feature,
            },
            {
                path: 'frontend',
                name: 'Frontend',
                meta: {
                    show: true,
                    label: '页面菜单',
                    icon: 'el-icon-menu',
                },
                component: Frontend,
            },
            {
                path: 'backend',
                name: 'Backend',
                meta: {
                    show: true,
                    label: '后台API',
                    icon: 'el-icon-picture',
                },
                component: Backend,
            },
        ],
    },
    // 组织管理
    {
        path: '/organization',
        name: 'Organization',
        meta: {
            show: true,
            label: '组织管理',
            icon: 'el-icon-office-building',
        },
        component: Organization,
        children: [
            {
                path: 'info',
                name: 'OrganizationInfo',
                meta: {
                    show: true,
                    label: '组织信息',
                    icon: 'el-icon-postcard',
                },
                component: OrganizationInfo,
            },
            {
                path: 'user',
                name: 'User',
                meta: {
                    show: true,
                    label: '用户管理',
                    icon: 'el-icon-user',
                },
                component: User,
            },
            {
                path: 'role',
                name: 'Role',
                meta: {
                    show: true,
                    label: '角色管理',
                    icon: 'el-icon-connection',
                },
                component: Role,
            },
            {
                path: 'object',
                name: 'OrganizationObject',
                meta: {
                    show: true,
                    label: '数据管理',
                    icon: 'el-icon-s-grid',
                },
                component: OrganizationObject,
            },

            {
                path: 'policy',
                name: 'OrganizationPolicy',
                meta: {
                    show: true,
                    label: '权限策略',
                    icon: 'el-icon-set-up',
                },
                component: OrganizationPolicy,
            },
        ],
    },
    // 404
    {
        path: '*',
        redirect: '/home',
    },
];

// 为了方便功能开发，默认关闭菜单权限控制
const menuRightsSwitch = true;

const filterAsideMenuList = async (routes) => {
    const menuList = [];
    await routes.reduce(async (memo, route) => {
        await memo;
        const { name, children } = route || {};
        const { show, rights } = route?.meta || {};
        if (children && children.length > 0) {
            const subMenu = await filterAsideMenuList(children);
            if (subMenu.length > 0) {
                route.children = subMenu;
                menuList.push(route);
            }
        } else {
            if (show) {
                if (rights === false || !menuRightsSwitch) {
                    menuList.push(route);
                } else {
                    const authResult = await Casbin.can({
                        object: `${name}#menu`,
                        action: 'read',
                    });
                    if (authResult) {
                        menuList.push(route);
                    }
                }
            }
        }
    }, undefined);
    return menuList;
};

export const getAsideMenuList = async () => {
    return await filterAsideMenuList(_.cloneDeep(routes));
};

const router = new VueRouter({
    routes,
});

router.beforeEach(async (to, from, next) => {
    // 页面访问权限校验
    const { name } = to || {};
    const { rights } = to?.meta || {};
    if (rights === false || !menuRightsSwitch) {
        next();
        return;
    }
    const authResult = await Casbin.can({
        object: `${name}#menu`,
        action: 'read',
    });
    if (authResult) {
        next();
    } else {
        next({
            name: 'Home',
        });
    }
});

export default router;
