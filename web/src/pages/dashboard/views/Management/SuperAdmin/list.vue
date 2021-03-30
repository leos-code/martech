<template>
    <el-table ref="adminTable" stripe :data="tableData">
        <el-table-column label="头像">
            <template slot-scope="scope">
                <img width="38px" height="38px" :src="scope.row.login_user_visible.avatar" class="avatar">
            </template>
        </el-table-column>
        <el-table-column prop="login_user_visible.nick_name" label="昵称" />
        <el-table-column prop="phone_number" label="手机号" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="updated_at_format" label="修改日期" />
        <el-table-column label="操作" min-width="200">
            <template slot-scope="scope">
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
                    <el-button type="danger" icon="el-icon-delete" size="mini" slot="reference">
                        删除
                    </el-button>
                </el-popover>
            </template>
        </el-table-column>
    </el-table>
</template>
<script>
import _ from 'lodash';
import { DateTimeFormat } from '@/utils';
import { SuperAdmin } from '@dashboard/api';
import { getVisibleLoginUser } from '@dashboard/views/Organization/User/utils';

export default {
    name: 'SuperAdminList',
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
        async handleDelete(index, row) {
            row.visible = false;
            const msgPrefix = '删除超级管理员';
            SuperAdmin.delete(row.id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        async getList() {
            const msgPrefix = '拉取超级管理员列表';
            SuperAdmin.get()
                .then(({ data }) => {
                    this.tableData = this.format(_.cloneDeep(data));
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        format(data) {
            data.forEach((item) => {
                item.updated_at_format = DateTimeFormat.formatByDate(item.updated_at, 'yyyy-mm-dd HH:MM');
                item.login_user_visible = getVisibleLoginUser(item.login_user);
            });
            return data;
        },
    },
};
</script>
