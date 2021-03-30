<template>
    <el-table ref="apiTable" stripe :data="tableData">
        <el-table-column prop="id" label="ApiId" sortable column-key="id" />
        <el-table-column prop="path" label="Path" />
        <el-table-column label="Method">
            <template slot-scope="scope">
                <el-tag :type="ApiMethods[scope.row.method] && ApiMethods[scope.row.method].type">
                    {{ scope.row.method }}
                </el-tag>
            </template>
        </el-table-column>
        <el-table-column prop="group" label="组" />
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
import { Backend } from '@dashboard/api';
import { ApiMethods } from './meta';

export default {
    name: 'BackendList',
    inject: ['reload'],
    data() {
        return {
            tableData: [],
            ApiMethods,
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
            const msgPrefix = '删除Api';
            row.visible = false;
            Backend.delete(row.id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        async getList() {
            const msgPrefix = '拉取Api列表';
            Backend.get()
                .then(({ data }) => {
                    this.tableData = this.format(_.clone(data));
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
