<template>
    <el-table ref="userTable" stripe :data="tableData">
        <el-table-column label="头像">
            <template slot-scope="scope">
                <img width="38px" height="38px" :src="scope.row.login_user_visible.avatar" class="avatar" />
            </template>
        </el-table-column>
        <el-table-column prop="login_user_visible.nick_name" label="昵称" />
        <el-table-column prop="phone_number" label="手机号" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column label="拥有角色">
            <template slot-scope="scope">
                <el-tag size="medium" v-for="role in scope.row.role" :key="role.id">
                    {{ role.name }}
                </el-tag>
            </template>
        </el-table-column>
        <el-table-column label="操作" min-width="200">
            <template slot-scope="scope">
                <el-button size="mini" icon="el-icon-edit" @click="handleEdit(scope.$index, scope.row)">
                    编辑
                </el-button>
                <el-popover v-model="scope.row.visible" placement="top" width="160">
                    <p>确定要解绑吗？</p>
                    <div style="text-align: right; margin: 0">
                        <el-button size="mini" type="text" @click="scope.row.visible = false">
                            取消
                        </el-button>
                        <el-button type="primary" size="mini" @click="handleDelete(scope.$index, scope.row)">
                            确定
                        </el-button>
                    </div>
                    <el-button slot="reference" type="danger" icon="el-icon-delete" size="mini">
                        解绑
                    </el-button>
                </el-popover>
            </template>
        </el-table-column>
    </el-table>
</template>
<script>
import _ from 'lodash';
import { DateTimeFormat } from '@/utils';
import { OrganizationUser } from '@dashboard/api';
import { getVisibleLoginUser } from '@dashboard/views/Organization/User/utils';

export default {
    name: 'UserList',
    inject: ['reload'],
    data() {
        return {
            tableData: [],
        };
    },
    mounted() {
        this.getList();
    },
    methods: {
        handleEdit(index, row) {
            this.$emit('edit', row);
        },
        handleDelete(index, row) {
            const msgPrefix = '解绑用户';
            row.visible = false;
            OrganizationUser.delete(row.id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        getList() {
            OrganizationUser.get()
                .then(({ data }) => {
                    this.tableData = this.format(_.clone(data));
                })
                .catch(({ msg }) => {
                    this.$message.error(`拉取用户列表失败 - ${msg}`);
                });
        },
        format(data) {
            data.map((item) => {
                item.updated_at_format = DateTimeFormat.formatByDate(item.updated_at, 'yyyy-mm-dd HH:MM');
                item.login_user_visible = getVisibleLoginUser(item.login_user);
            });
            return data;
        },
    },
};
</script>
