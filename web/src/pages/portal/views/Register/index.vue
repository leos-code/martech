<template>
    <fullscreen>
        <el-container class="register-container">
            <el-tabs class="register-tab" type="border-card" :stretch="true">
                <el-tab-pane class="register-tab-pane" label="用户注册">
                    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
                        <el-form-item label="昵称">
                            <el-input v-model="nickname" disabled />
                        </el-form-item>
                        <el-form-item label="头像">
                            <el-avatar :src="avatar" />
                        </el-form-item>
                        <el-form-item label="手机号" prop="phone_number">
                            <el-input v-model="form.phone_number" autocomplete="off" />
                        </el-form-item>
                        <el-form-item label="邮箱" prop="email">
                            <el-input v-model="form.email" autocomplete="off" @keyup.enter.native="commit" />
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" @click="commit">
                                注册
                            </el-button>
                        </el-form-item>
                    </el-form>
                </el-tab-pane>
            </el-tabs>
        </el-container>
    </fullscreen>
</template>
<script>
import Fullscreen from '@portal/views/Layout/fullscreen';
import Register from '@portal/api/register';

export default {
    name: 'Register',
    components: {
        Fullscreen,
    },
    data() {
        return {
            loginUser: {},
            formRef: 'registerForm',
            form: {
                phone_number: '',
                email: '',
            },
            nickname: '',
            avatar: '',
            rules: {
                phone_number: [
                    {
                        required: true,
                        message: '请输入手机号',
                    },
                ],
                email: [
                    {
                        required: true,
                        message: '请输入邮箱',
                    },
                ],
            },
        };
    },
    mounted() {
        this.init();
    },
    methods: {
        async init() {
            await this.getLoginUser();
            const { avatar, nick_name } = this.loginUser || {};
            this.avatar = avatar;
            this.nickname = nick_name;
        },
        getLoginUser() {
            return new Promise((resolve, reject) => {
                const msgPrefix = '拉取登录信息';
                Register.get()
                    .then(({ data }) => {
                        this.loginUser = data;
                        resolve();
                    })
                    .catch(({ msg }) => {
                        this.$message.error(`${msgPrefix}失败 - ${msg}`);
                        reject();
                    });
            });
        },
        commit() {
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    const msgPrefix = '注册';
                    Register.create(this.form)
                        .then(() => {
                            this.$message.success(`${msgPrefix}成功, 稍后将自动跳转管理台`);
                            setTimeout(() => {
                                window.location.href = process.env.VUE_APP_PAGE_DASHBOARD;
                            }, 1500);
                        })
                        .catch(({ msg }) => {
                            this.$message.error(`${msgPrefix}失败 - ${msg}`);
                        });
                } else {
                    return false;
                }
            });
        },
    },
};
</script>
<style lang="scss">
.register-container {
    height: 100%;
    justify-content: center;
    .register-tab {
        width: 710px;
        height: 600px;
        margin-top: 60px;
        .el-tabs__item {
            height: 60px;
            line-height: 60px;
            font-size: 18px;
        }
        .register-tab-pane {
            height: 520px;
        }
    }
}
</style>
