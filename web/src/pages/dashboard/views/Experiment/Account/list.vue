<template>
    <el-table ref="experimentParameterTable" stripe :data="tableData" lazy>
        <el-table-column type="expand">
            <template slot-scope="prop">
                <rta-account-detail :account="prop.row" />
            </template>
        </el-table-column>
        <el-table-column prop="id" label="账户ID" sortable column-key="id" />
        <el-table-column prop="rta_id" label="RTA ID" />
        <el-table-column prop="token" label="Token" />
        <el-table-column prop="description" label="账户描述" />
        <el-table-column prop="updated_at_format" label="修改日期" />
        <el-table-column label="操作" min-width="200">
            <template slot-scope="scope">
                <el-button size="mini" icon="el-icon-edit" @click="handleEdit(scope.$index, scope.row)">
                    编辑
                </el-button>
                <el-popover v-model="scope.row.visible" placement="top" width="160">
                    <p>确定要废弃吗？</p>
                    <div style="text-align: right; margin: 0">
                        <el-button size="mini" type="text" @click="scope.row.visible = false">
                            取消
                        </el-button>
                        <el-button type="primary" size="mini" @click="handleDelete(scope.$index, scope.row)">
                            确定
                        </el-button>
                    </div>
                    <el-button slot="reference" type="danger" icon="el-icon-delete" size="mini">
                        废弃
                    </el-button>
                </el-popover>
                <el-button size="mini" type="primary" icon="el-icon-edit" @click="handleSync(scope.$index, scope.row)">
                    刷新
                </el-button>
            </template>
        </el-table-column>
    </el-table>
</template>
<script>
import _ from 'lodash';
import { DateTimeFormat } from '@/utils';
import { ExperimentAccount } from '@dashboard/api';
import RtaAccountDetail from '@dashboard/views/Experiment/Account/detail';

export default {
    name: 'ExperimentAccountList',
    inject: ['reload'],
    components: {
        RtaAccountDetail,
    },
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
            const msgPrefix = '废弃账户';
            row.visible = false;
            ExperimentAccount.delete(row.id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        getList() {
            const msgPrefix = '拉取账户列表';
            ExperimentAccount.getList()
                .then(({ data: accountList }) => {
                    this.tableData = this.format(_.clone(accountList));
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
        handleSync(index, row) {
            const msgPrefix = '同步账户信息';
            ExperimentAccount.sync(row)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
    },
};
</script>
