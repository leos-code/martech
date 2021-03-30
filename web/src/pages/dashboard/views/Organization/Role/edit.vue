<template>
    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
        <el-form-item label="父角色" label-width="160px">
            <el-cascader
                :options="roleList"
                :props="parentRoleProps"
                clearable
                placeholder="不选默认为根角色"
                v-model="form.parent_id"
                :show-all-levels="false"
            />
        </el-form-item>
        <el-form-item label="角色名称" label-width="160px" prop="name">
            <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="绑定用户" label-width="160px">
            <el-select v-model="form.userIds" multiple>
                <el-option
                    v-for="user in userList"
                    :key="user.id"
                    :label="user.login_user_visible && user.login_user_visible.nick_name"
                    :value="user.id"
                />
            </el-select>
        </el-form-item>
        <el-form-item label="绑定数据" label-width="160px">
            <el-tree
                :ref="policySelectRef"
                :data="policyList"
                show-checkbox
                node-key="id"
                default-expand-all
                :props="policyTreeProps"
                check-on-click-node
                :expand-on-click-node="false"
                :check-strictly="true"
            >
                <span class="custom-tree-node" slot-scope="{ data }">
                    <span>{{ data && data.name }}</span>
                    <span>
                        <el-checkbox v-if="data && !data.disableRead" v-model="data.read">读</el-checkbox>
                        <el-checkbox v-if="data && !data.disableWrite" v-model="data.write">写</el-checkbox>
                    </span>
                </span>
            </el-tree>
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
import { OrganizationUser, Role, ObjectApi, Policy } from '@dashboard/api';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import Utils from '@/utils';
import { getVisibleLoginUser } from '@dashboard/views/Organization/User/utils';

const DEFAULT_FORM = {
    name: '',
    parent_id: 0,
    user: [],
};

export default {
    name: 'RoleEdit',
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
        parentValue: {
            type: Object,
            default: () => {
                return {};
            },
        },
    },
    data() {
        return {
            formRef: 'roleEditForm',
            form: _.cloneDeep(DEFAULT_FORM),
            rules: {
                name: [
                    {
                        required: true,
                        message: '请输入角色名称',
                    },
                ],
            },
            roleList: [],
            parentRoleProps: {
                value: 'id',
                label: 'name',
                children: 'children',
                checkStrictly: true,
            },
            userList: [],
            DIALOG_TYPE,
            objectList: [],
            policyList: [],
            policySelectRef: 'policySelectRef',
            policyTreeProps: {
                children: 'children',
                label: 'name',
            },
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

            await this.getRoleList();

            // 初始化默认选中的父节点
            if (this.type === DIALOG_TYPE.ADD && Utils.isPositiveInteger(this.parentValue?.id)) {
                this.form.parent_id = this.parentValue.id;
            }

            await this.getUserList();
            this.setDefaultSelectedUsers();

            await this.getObjectList();
            this.setDefaultSelectedPolicies();

            await this.getDefaultRoleObject();
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
        getUserList() {
            const msgPrefix = '拉取用户列表';
            return new Promise((resolve, reject) => {
                if (Utils.isNonEmptyArray(this.userList)) {
                    resolve();
                    return;
                }
                OrganizationUser.get()
                    .then(({ data }) => {
                        this.userList = this.formatUser(_.cloneDeep(data));
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
        setDefaultSelectedUsers() {
            const { user: users } = this.form;
            if (Utils.isNonEmptyArray(users)) {
                this.$set(
                    this.form,
                    'userIds',
                    users.map((user) => user.id),
                );
            }
        },
        getDefaultRoleObject() {
            const msgPrefix = '拉取角色Object';
            return new Promise((resolve, reject) => {
                ObjectApi.get({
                    type: 'role',
                })
                    .then(({ data }) => {
                        if (Utils.isNonEmptyArray(data)) {
                            const { id } = data[0] || {};
                            this.form.object = `object_${id}`;
                        }
                        resolve();
                    })
                    .catch((error) => {
                        console.error(error);
                        const { msg } = error || {};
                        this.$message.error(`${msgPrefix}失败 - ${msg || '未知'}`);
                        reject();
                    });
            });
        },
        getObjectList() {
            const msgPrefix = '拉取数据列表';
            return new Promise((resolve, reject) => {
                if (Utils.isNonEmptyArray(this.objectList)) {
                    resolve();
                    return;
                }
                ObjectApi.get()
                    .then(({ data }) => {
                        this.objectList = this.formatObject(data);
                        resolve();
                    })
                    .catch((error) => {
                        console.error(error);
                        const { msg } = error || {};
                        this.$message.error(`${msgPrefix}失败 - ${msg || '未知'}`);
                        reject();
                    });
            });
        },
        setDefaultSelectedPolicies() {
            this.policyList = _.cloneDeep(this.objectList);
            const { policy: policies } = this.form;
            let defaultSelectedPolicyIds = [];
            if (Utils.isNonEmptyArray(policies)) {
                policies.forEach((policy) => {
                    const { object, read = false, write = false } = policy || {};
                    const { id } = object || {};
                    const selectedPolicyIndex = _.findIndex(this.policyList, {
                        id,
                    });
                    if (selectedPolicyIndex !== -1) {
                        defaultSelectedPolicyIds.push(id);
                        this.$set(this.policyList[selectedPolicyIndex], 'read', read);
                        this.$set(this.policyList[selectedPolicyIndex], 'write', write);
                    }
                });
            }

            this.policyList = Utils.arrayToTree(this.policyList, {
                parentIdNullVal: 0,
            });

            this.$refs[this.policySelectRef].setCheckedKeys(defaultSelectedPolicyIds);
        },
        formatObject(data) {
            data.map((object) => {
                const { type } = object || {};
                // 功能数据不支持 写 权限
                if (type === 'feature') {
                    object.disableWrite = true;
                }
            });
            return data;
        },
        formatUser(data) {
            data.map((item) => {
                item.login_user_visible = getVisibleLoginUser(item.login_user);
            });
            return data;
        },
        reset() {
            this.form = _.cloneDeep(DEFAULT_FORM);
            this.$refs[this.formRef].resetFields();
        },
        close() {
            this.$emit('close');
            this.reset();
        },
        getSelectedPolicy() {
            const selectedPolicies = this.$refs[this.policySelectRef].getCheckedNodes();
            const resultPolicies = [];
            selectedPolicies.forEach((policy) => {
                const { id, read = false, write = false } = policy || {};
                const originObject = _.find(this.objectList, {
                    id,
                });
                if (typeof originObject !== 'undefined') {
                    resultPolicies.push({
                        object: originObject,
                        read,
                        write,
                    });
                }
            });
            return resultPolicies;
        },
        async commit() {
            const msgPrefix = this.type === DIALOG_TYPE.ADD ? '新建角色' : '修改角色';

            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    // 将用户id列表转为用户结构列表
                    if (Utils.isNonEmptyArray(this.form.userIds)) {
                        this.form.user = _.filter(this.userList, (o) => {
                            return this.form.userIds.indexOf(o.id) > -1;
                        });
                    } else {
                        this.form.user = [];
                    }
                    // 将parent_id列表转为单数字
                    if (Utils.isArray(this.form.parent_id)) {
                        const length = this.form.parent_id.length;
                        if (length > 0) {
                            this.form.parent_id = this.form.parent_id[length - 1];
                        } else {
                            this.form.parent_id = 0;
                        }
                    }

                    // 获取选中的权限列表
                    this.form.policy = this.getSelectedPolicy();

                    Role.edit(this.form)
                        .then(({ data }) => {
                            this.form = _.merge({}, this.form, data);
                            Role.editUser(this.form).then(() => {
                                Policy.edit(this.form).then(() => {
                                    this.$message.success(`${msgPrefix}成功`);
                                    this.reload();
                                });
                            });
                        })
                        .catch((error) => {
                            const { msg } = error || {};
                            console.error(error);
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
.custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 13px;
    padding-right: 8px;
    .el-checkbox__label {
        font-size: 13px;
    }
}
</style>
