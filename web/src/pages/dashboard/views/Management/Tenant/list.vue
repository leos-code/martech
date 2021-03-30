<template>
    <el-table ref="tenantTable" stripe :data="tableData">
        <el-table-column prop="id" label="租户Id" sortable column-key="id" />
        <el-table-column prop="name" label="租户名称" />
        <el-table-column prop="updated_at_format" label="修改日期" />
        <el-table-column label="操作" min-width="200">
            <template slot-scope="scope">
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
                    <el-button slot="reference" type="danger" icon="el-icon-delete" size="mini">
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
import { Tenant } from '@dashboard/api';

export default {
    name: 'TenantList',
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
        async handleDelete(index, row) {
            const msgPrefix = '删除租户';
            row.visible = false;
            Tenant.delete(row.id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        async getList() {
            const msgPrefix = '拉取租户列表';
            Tenant.get()
                .then(({ data }) => {
                    this.tableData = this.format(_.clone(data));
                    this.$store.dispatch('user/getUserInfo');
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        format(data) {
            data.map((item) => {
                item.updated_at_format = DateTimeFormat.formatByDate(item.updated_at, 'yyyy-mm-dd HH:MM');
            });
            return data;
        },
    },
};
</script>
