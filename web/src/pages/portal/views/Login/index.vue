<template>
    <fullscreen>
        <el-container class="login-container">
            <el-tabs class="login-tab" type="border-card" :stretch="true">
                <template v-for="oauth in oauthList">
                    <el-tab-pane :key="oauth.type" :label="oauth.label" class="login-tab-pane">
                        <iframe frameborder="0" allowtransparency="yes" :src="oauth.url" width="100%" height="100%" />
                    </el-tab-pane>
                </template>
            </el-tabs>
        </el-container>
    </fullscreen>
</template>
<script>
import Fullscreen from '@portal/views/Layout/fullscreen';
import Login from '@portal/api/login';
import { LOGIN_DIALOG_LABEL } from './meta';
import Utils from '@/utils';

export default {
    name: 'Login',
    components: {
        Fullscreen,
    },
    data() {
        return {
            oauthList: [],
        };
    },
    methods: {
        init() {
            this.getOauthList();
            this.setGlobalCallback();
        },
        getOauthList() {
            const msgPrefix = '获取登录链接';
            Login.oauth()
                .then(({ data }) => {
                    this.oauthList = data;
                    if (Utils.isNonEmptyArray(this.oauthList)) {
                        this.oauthList.forEach((oauth) => {
                            const { type } = oauth || {};
                            oauth.label = LOGIN_DIALOG_LABEL[type].label || '账户登录';
                        });
                    }
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        setGlobalCallback() {
            if (!Utils.isFunction(window.loginCallback)) {
                window.loginCallback = (options) => {
                    const { redirect_type } = options || {};
                    if (redirect_type === 'dashboard') {
                        window.location.href = process.env.VUE_APP_PAGE_DASHBOARD;
                    } else if (redirect_type === 'register') {
                        this.$router.push({
                            name: 'Register',
                        });
                    }
                };
            }
        },
    },
    mounted() {
        this.init();
    },
};
</script>
<style lang="scss">
.login-container {
    height: 100%;
    justify-content: center;
    .login-tab {
        width: 1000px;
        height: 780px;
        margin-top: 60px;
        .el-tabs__item {
            height: 60px;
            line-height: 60px;
            font-size: 18px;
        }
        .login-tab-pane {
            height: 680px;
        }
    }
}
</style>
