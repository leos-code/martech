<template>
    <el-table ref="strategyTable" stripe :data="tableData">
        <el-table-column prop="id" label="策略Id" sortable column-key="id" />
        <el-table-column prop="name" label="策略名称" />
        <el-table-column prop="platform" label="平台" />
        <el-table-column prop="targeting_id" label="定向ID" />
        <el-table-column prop="strategy.strategy_id" label="平台RTA Id" />
        <el-table-column prop="strategy.advertiser_id" label="广告主Id" />
        <el-table-column prop="strategy.campaign_id" label="推广计划Id" />
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
import { BindStrategy } from '@dashboard/api';

export default {
    name: 'BindStrategyList',
    inject: ['reload'],
    data() {
        return {
            tableData: [],
        };
    },
    mounted() {
        console.log('BindStrategy List mounted');
        this.getList();
    },
    created() {
        console.log('BindStrategy List created');
    },
    methods: {
        handleEdit(index, row) {
            this.$emit('edit', row);
        },
        async handleDelete(index, row) {
            row.visible = false;
            const response = await BindStrategy.delete(row.id);
            const { data } = response || {};
            const { code, msg } = data || {};
            if (code === 0) {
                this.$message({
                    message: '删除成功',
                    type: 'success',
                });
                this.reload();
            } else {
                this.$message({
                    message: `删除失败 - ${msg}`,
                    type: 'danger',
                });
            }
        },
        async getList() {
            const response = await BindStrategy.getList();
            const { data } = response || {};
            const { code, msg, strategy } = data || {};
            if (code === 0) {
                this.tableData = this.format(_.clone(strategy));
            } else {
                this.$message({
                    message: msg,
                    type: 'error',
                });
            }
        },
        format(data) {
            data.map((item) => {
                const date = new Date(item.updated_at);
                item.updated_at_format = date.toLocaleString();
            });
            return data;
        },
    },
};
</script>
