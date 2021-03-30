<template>
    <el-table ref="experimentParameterTable" stripe :data="tableData">
        <el-table-column prop="id" label="参数Id" sortable column-key="id" />
        <el-table-column prop="name" label="参数名称" />
        <el-table-column prop="type_name" label="类型" />
        <el-table-column prop="description" label="参数描述" />
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
            </template>
        </el-table-column>
    </el-table>
</template>
<script>
import { mapGetters } from 'vuex';
import _ from 'lodash';
import { DateTimeFormat } from '@/utils';
import { PARAMETER_TYPE } from './meta';
import { ExperimentParameter } from '@dashboard/api';

export default {
    name: 'ExperimentParameterList',
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
            const msgPrefix = '废弃参数';
            row.visible = false;
            ExperimentParameter.delete(row.id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        async getList() {
            ExperimentParameter.getList()
                .then(({ data }) => {
                    this.tableData = this.format(_.clone(data));
                })
                .catch(({ msg }) => {
                    this.$message.error(`拉取参数失败 - ${msg}`);
                });
        },
        format(data) {
            data.map((item) => {
                item.updated_at_format = DateTimeFormat.formatByDate(item.updated_at, 'yyyy-mm-dd HH:MM');
                const parameterType = _.find(PARAMETER_TYPE, {
                    key: item.type,
                });
                item.type_name = (parameterType && parameterType.name) || '未知';
            });
            return data;
        },
    },
};
</script>
