<template>
    <el-form :ref="formRef" :model="form" label-width="80px">
        <template v-if="type === DIALOG_TYPE.ADD">
            <user-search :ref="userSearchRef" label-width="160px" @select="handleUserSelected" />
        </template>
        <template v-else>
            <el-form-item label="昵称" label-width="160px">
                <img
                    width="80px"
                    height="80px"
                    :src="form && form.login_user_visible && form.login_user_visible.avatar"
                    class="avatar"
                />
                <el-input :value="form && form.login_user_visible && form.login_user_visible.nick_name" disabled />
            </el-form-item>
            <el-form-item label="手机号" label-width="160px">
                <el-input :value="form.phone_number" disabled />
            </el-form-item>
            <el-form-item label="邮箱" label-width="160px">
                <el-input :value="form.email" disabled />
            </el-form-item>
        </template>
        <el-form-item label="选择角色">
            <el-tree
                :ref="roleSelectRef"
                :data="roleList"
                show-checkbox
                node-key="id"
                default-expand-all
                :props="roleTreeProps"
                check-on-click-node
                :default-checked-keys="defaultSelectedRoles"
                :expand-on-click-node="false"
                :check-strictly="true"
            />
        </el-form-item>
        <el-form-item>
            <el-button @click="close">
                取 消
            </el-button>
            <el-button type="primary" @click="commit">
                {{ commitText }}
            </el-button>
        </el-form-item>
    </el-form>
</template>
<script>
import _ from 'lodash';
import { OrganizationUser, Role } from '@dashboard/api';
import UserSearch from './search';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import Utils from '@/utils';

const DEFAULT_FORM = {
    name: '',
};

export default {
    name: 'UserEdit',
    inject: ['reload'],
    props: {
        type: {
            type: String,
            default: DIALOG_TYPE.ADD,
        },
        commitText: String,
        defaultValue: {
            type: Object,
            default: () => {
                return {};
            },
        },
    },
    components: {
        UserSearch,
    },
    data() {
        return {
            formRef: 'userEditForm',
            userSearchRef: 'userSearchRef',
            roleSelectRef: 'roleSelectRef',
            form: _.clone(DEFAULT_FORM),
            rules: {
                name: [
                    {
                        required: true,
                        message: '请输入账户名',
                    },
                ],
            },
            roleList: [],
            roleTreeProps: {
                children: 'children',
                label: 'name',
            },
            defaultSelectedRoles: [],
            DIALOG_TYPE,
        };
    },
    watch: {
        defaultValue() {
            this.init();
        },
    },
    mounted() {
        this.init();
    },
    methods: {
        async init() {
            if (typeof this.defaultValue === 'undefined') {
                this.reset();
            } else {
                this.form = _.assign({}, DEFAULT_FORM, this.defaultValue);
            }
            this.defaultSelectedRoles = [];
            await this.getRoleList();
            this.setDefaultSelectedRoles();
        },
        getRoleList() {
            const msgPrefix = '拉取角色列表';
            return new Promise((resolve, reject) => {
                if (Utils.isNonEmptyArray(this.roleList)) {
                    resolve();
                    return;
                }
                Role.get()
                    .then(({ data }) => {
                        this.roleList = Utils.arrayToTree(_.cloneDeep(data), {
                            parentIdNullVal: 0,
                        });
                        resolve();
                    })
                    .catch((error) => {
                        const { msg } = error || {};
                        this.$message.error(`${msgPrefix}失败 - ${msg || '未知'}`);
                        console.error(error);
                        reject();
                    });
            });
        },
        setDefaultSelectedRoles() {
            const { role = [] } = this.form || {};
            this.$refs[this.roleSelectRef].setCheckedNodes(role);
        },
        reset() {
            this.$refs[this.userSearchRef]?.resetFields();
            this.defaultSelectedRoles = [];
        },
        close() {
            this.$emit('close');
            this.reset();
        },
        async commit() {
            const msgPrefix = this.type === DIALOG_TYPE.ADD ? '新建用户' : '修改用户';
            const { id } = this.form || {};
            // 用户选中状态校验
            if (!Utils.isPositiveInteger(id)) {
                this.$message.warning(`${msgPrefix} - 请先选中一个用户`);
                return false;
            }
            // 获取选中的角色列表
            const selectedRoles = this.$refs[this.roleSelectRef].getCheckedNodes();
            this.form.role = selectedRoles;

            OrganizationUser.edit(this.form)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        handleUserSelected(user) {
            const { id } = user || {};
            if (Utils.isPositiveInteger(id)) {
                this.form = _.merge(this.form, user);
            }
        },
    },
};
</script>
