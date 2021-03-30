<template>
    <el-aside class="layout-aside">
        <div class="title">
            <img src="~@/assets/logo.png" alt="" class="logo-img">
            <h2 class="title-text">
                User Growth
            </h2>
        </div>
        <el-menu
            unique-opened
            background-color="#191a23"
            text-color="rgb(191, 203, 217)"
            active-text-color="#fff"
            :default-active="defaultActive"
            :collapse-transition="true"
            @select="select"
        >
            <template v-for="menu in menuList">
                <template v-if="menu.children && menu.children.length">
                    <el-submenu :key="menu.name" :index="menu.name">
                        <template slot="title">
                            <i :class="menu.meta.icon" />
                            <span>{{ menu.meta.label }}</span>
                        </template>
                        <template v-for="subMenu in menu.children">
                            <el-menu-item :key="subMenu.name" :index="subMenu.name">
                                <i :class="subMenu.meta.icon" />
                                <span>{{ subMenu.meta.label }}</span>
                            </el-menu-item>
                        </template>
                    </el-submenu>
                </template>
                <template v-else>
                    <el-menu-item :key="menu.name" :index="menu.name">
                        <i :class="menu.meta.icon" />
                        <span>{{ menu.meta.label }}</span>
                    </el-menu-item>
                </template>
            </template>
        </el-menu>
    </el-aside>
</template>
<script>
import _ from 'lodash';
import { getAsideMenuList } from '@dashboard/router';

export default {
    name: 'UgAside',
    data() {
        return {
            menuList: [],
            defaultActive: '',
        };
    },
    watch: {
        $route() {
            this.defaultActive = this.getActiveMenu();
        },
    },
    async mounted() {
        this.menuList = await getAsideMenuList();
        console.log(this.menuList);
    },
    methods: {
        select(name) {
            const { name: currentName } = this.$route || {};
            if (currentName === name) {
                return;
            }
            this.$router.push({
                name,
            });
        },
        getActiveMenu() {
            return this.$route?.name;
        },
    },
};
</script>
