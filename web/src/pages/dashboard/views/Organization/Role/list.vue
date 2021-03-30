<template>
    <el-table
        :ref="roleTable"
        stripe
        :data="tableData"
        row-key="id"
        default-expand-all
        border
        :tree-props="roleTreeProps"
    >
        <el-table-column prop="name" label="角色名称" />
        <el-table-column label="绑定用户">
            <template slot-scope="scope">
                <el-popover v-for="user in scope.row.user" :key="user.id" trigger="hover" placement="top">
                    <p>昵称: {{ user && user.login_user_visible && user.login_user_visible.nick_name }}</p>
                    <img
                        slot="reference"
                        width="38px"
                        height="38px"
                        :src="user && user.login_user_visible && user.login_user_visible.avatar"
                        class="avatar"
                    />
                </el-popover>
            </template>
        </el-table-column>
        <el-table-column label="绑定数据">
            <template slot-scope="scope">
                <el-tag size="medium" v-for="(policy, index) in scope.row.policy" :key="index">
                    {{ policy.object && policy.object.name }}
                </el-tag>
            </template>
        </el-table-column>
        <el-table-column label="操作" min-width="200">
            <template slot-scope="scope">
                <el-button
                    size="mini"
                    icon="el-icon-plus"
                    type="primary"
                    @click="handleCreate(scope.$index, scope.row)"
                >
                    创建子角色
                </el-button>
                <el-button size="mini" icon="el-icon-edit" @click="handleEdit(scope.$index, scope.row)">
                    编辑
                </el-button>
                <el-popover v-model="scope.row.visible" placement="top" width="160">
                    <p>确定要删除吗？</p>
                    <div style="text-align: right; margin: 0">
                        <el-button size="mini" type="text" @click="scope.row.visible = false">
                            取消
                        </el-button>
                        <el-button type="primary" size="mini" @click="handleDelete(scope.$index, scope.row)">
                            确定
                        </el-button>
                    </div>
                    <el-button slot="reference" type="danger" size="mini" icon="el-icon-delete">
                        删除
                    </el-button>
                </el-popover>
            </template>
        </el-table-column>
    </el-table>
</template>
<script>
import _ from 'lodash';
import { Role, Policy } from '@dashboard/api';
import { getVisibleLoginUser } from '@dashboard/views/Organization/User/utils';
import Utils from '@/utils';

export default {
    name: 'RoleList',
    inject: ['reload'],
    data() {
        return {
            roleTable: 'roleTable',
            tableData: [],
            roleTreeProps: {
                children: 'children',
            },
        };
    },
    mounted() {
        this.getList();
    },
    methods: {
        handleCreate(index, row) {
            this.$emit('create', row);
        },
        handleEdit(index, row) {
            this.$emit('edit', row);
        },
        handleDelete(index, row) {
            const msgPrefix = '删除角色';
            row.visible = false;
            Role.delete(row.id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        getList() {
            Role.get()
                .then(({ data: roleList }) => {
                    Policy.get().then(({ data: policyList }) => {
                        roleList.forEach((role) => {
                            const object = _.find(policyList, {
                                id: role?.id,
                            });
                            if (typeof object !== 'undefined') {
                                role.policy = object.policy;
                            }
                        });
                        this.tableData = this.format(roleList);
                    });
                })
                .catch((error) => {
                    const { msg } = error || {};
                    this.$message.error(`拉取角色列表失败 - ${msg}`);
                    console.error(error);
                });
        },
        format(data) {
            data.map((item) => {
                if (Utils.isNonEmptyArray(item.user)) {
                    item.user.forEach((user) => {
                        user.login_user_visible = getVisibleLoginUser(user.login_user);
                    });
                }
            });
            data = Utils.arrayToTree(data, {
                parentIdNullVal: 0,
            });
            return data;
        },
    },
};
</script>
