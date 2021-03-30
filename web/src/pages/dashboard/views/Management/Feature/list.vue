<template>
    <el-table ref="menuTable" stripe :data="tableData">
        <el-table-column prop="id" label="功能Id" sortable column-key="id" />
        <el-table-column prop="name" label="功能名称" />
        <el-table-column prop="description" label="描述" />
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
import { Feature } from '@dashboard/api';

export default {
    name: 'FeatureList',
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
            const msgPrefix = '删除功能';
            row.visible = false;
            Feature.delete(row.id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        async getList() {
            const msgPrefix = '拉取功能列表';
            Feature.get()
                .then(({ data }) => {
                    this.tableData = this.format(_.cloneDeep(data));
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        format(data) {
            data = _.filter(data, (o) => {
                return typeof o.parent_id !== 'undefined';
            });
            data.map((item) => {
                item.updated_at_format = DateTimeFormat.formatByDate(item.updated_at, 'yyyy-mm-dd HH:MM');
            });
            return data;
        },
    },
};
</script>
