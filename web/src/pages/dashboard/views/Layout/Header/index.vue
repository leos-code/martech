<template>
    <div>
        <el-header class="layout-header">
            <div class="tenant-info" @click="handleTenantSwitch">
                <span v-if="currentTenant && currentTenant.name">组织名称：{{ currentTenant.name }}</span>
                <span v-else>您当前未选择组织，请点击选择</span>
            </div>
            <div class="staff-info">
                <template v-if="loginUser.nick_name">
                    <el-dropdown @command="handleCommand">
                        <img
                            alt="头像"
                            width="38"
                            height="38"
                            style="border-radius: 50%; margin-top:10px;"
                            :src="loginUser.avatar"
                        >
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item command="info">
                                {{ loginUser.nick_name }}
                            </el-dropdown-item>
                            <el-dropdown-item command="logout" divided>
                                退出登录
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </template>
            </div>
        </el-header>
    </div>
</template>
<script>
import { mapGetters } from 'vuex';
import Login from '@portal/api/login';
import Casbin from '@dashboard/Casbin';

export default {
    name: 'UgHeader',
    computed: {
        ...mapGetters('user', ['loginUser', 'currentTenant']),
    },
    methods: {
        handleCommand(command) {
            switch (command) {
                case 'info':
                    console.log('跳转个人信息页');
                    break;
                case 'logout': {
                    console.log('退出登录');
                    const msgPrefix = '退出登录';
                    Login.logout()
                        .then(async () => {
                            await Casbin.clear();
                            window.location.href = process.env.VUE_APP_PAGE_PORTAL;
                        })
                        .catch(({ msg }) => {
                            this.$message.error(`${msgPrefix}失败 - ${msg}`);
                        });
                    break;
                }
                default:
                    break;
            }
        },
        handleTenantSwitch() {
            this.$store.dispatch('user/activeTenantSwitch');
        },
    },
};
</script>
<style lang="scss">
.layout-header {
    background-color: #ffffff;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    text-align: right;
    .tenant-info {
        cursor: pointer;
        display: inline-block;
        height: 30px;
        width: 200px;
        margin: 10px 20px 0 0;
        padding: 5px 10px;
        vertical-align: top;
        text-align: left;
        border: 1px solid #bbbbbb;
        color: #555555;
        border-radius: 15px;
        line-height: 30px;
        font-size: 14px;
    }
    .staff-info {
        cursor: pointer;
        display: inline-block;
        margin-right: 20px;
    }
}
</style>
