<template>
    <el-container class="layout-container">
        <ug-aside />
        <el-main class="main-container">
            <ug-header />
            <router-view v-if="isRouterAlive" class="main-box" />
        </el-main>
        <tenant-switch />
    </el-container>
</template>
<script>
import '@/style/layout.scss';

import UgHeader from '@dashboard/views/Layout/Header/index';
import UgAside from '@dashboard/views/Layout/Aside/index';
import { TenantSwitch } from '@dashboard/views/Management';

export default {
    name: 'UgLayout',
    components: {
        UgHeader,
        UgAside,
        TenantSwitch,
    },
    provide() {
        return {
            reload: this.reload,
        };
    },
    data() {
        return {
            isRouterAlive: true,
        };
    },
    methods: {
        reload() {
            console.log('reload call');
            this.isRouterAlive = false;
            this.$nextTick(() => {
                this.isRouterAlive = true;
            });
        },
    },
};
</script>
