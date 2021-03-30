<template>
    <el-table
        :ref="objectTable"
        stripe
        :data="tableData"
        row-key="id"
        border
        :tree-props="objectTreeProps"
        default-expand-all
    >
        <el-table-column prop="name" label="名称" />
        <el-table-column label="操作" min-width="200">
            <template slot-scope="scope">
                <el-button
                    size="mini"
                    icon="el-icon-plus"
                    type="primary"
                    @click="handleCreate(scope.$index, scope.row)"
                >
                    创建{{ objectType.label }}
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
import { ObjectApi } from '@dashboard/api';
import Utils from '@/utils';

export default {
    name: 'ObjectList',
    inject: ['reload'],
    props: {
        objectType: Object,
    },
    data() {
        return {
            objectTable: 'objectTable',
            tableData: [],
            objectTreeProps: {
                children: 'children',
            },
        };
    },
    mounted() {
        this.getList();
    },
    methods: {
        handleCreate(index, row) {
            this.$emit('create', this.objectType, row);
        },
        handleEdit(index, row) {
            this.$emit('edit', this.objectType, row);
        },
        handleDelete(index, row) {
            const msgPrefix = `删除${this.objectType?.label}`;
            row.visible = false;
            ObjectApi.delete(row.id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
        getList() {
            const { type, label } = this.objectType || {};
            const msgPrefix = `拉取${label}信息`;
            ObjectApi.get({
                type,
            })
                .then(({ data }) => {
                    this.tableData = this.format(_.clone(data));
                })
                .catch((error) => {
                    console.error(error);
                    const { msg } = error || {};
                    this.$message.error(`${msgPrefix}失败 - ${msg || '未知'}`);
                });
        },
        format(data) {
            data = Utils.arrayToTree(data, {
                parentIdNullVal: 0,
            });
            return data;
        },
    },
};
</script>
