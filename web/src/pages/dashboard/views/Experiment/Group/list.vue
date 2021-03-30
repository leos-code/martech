<template>
    <el-table ref="experimentGroupTable" stripe :data="tableData" lazy>
        <el-table-column prop="id" label="实验组ID" sortable column-key="id" />
        <el-table-column prop="name" label="实验组名称" column-key="id" />
        <el-table-column prop="updated_at_format" label="修改时间" column-key="id" />
        <el-table-column label="操作" min-width="200">
            <template slot-scope="scope">
                <el-button size="mini" icon="el-icon-edit" type="primary" @click="handleRadar(scope.$index, scope.row)">
                    效果
                </el-button>
                <el-button size="mini" icon="el-icon-edit" @click="handleDraft(scope.$index, scope.row)">
                    管理
                </el-button>
            </template>
        </el-table-column>
    </el-table>
</template>
<script>
import _ from 'lodash';
import { DateTimeFormat } from '@/utils';
import { ExperimentGroup } from '@dashboard/api';

export default {
    name: 'ExperimentGroupList',
    data() {
        return {
            tableData: [],
        };
    },
    mounted() {
        this.getList();
    },
    methods: {
        getList() {
            ExperimentGroup.getList()
                .then(({ data }) => {
                    this.tableData = this.format(_.clone(data));
                })
                .catch(({ msg }) => {
                    this.$message.error(`拉取实验组失败 - ${msg}`);
                });
        },
        format(groups) {
            groups.map((item) => {
                item.updated_at_format = DateTimeFormat.formatByDate(item.updated_at, 'yyyy-mm-dd HH:MM');
            });
            return groups;
        },
        handleDraft(index, row) {
            this.$router.push({
                name: 'ExperimentGroupDetail',
                params: {
                    id: row.id,
                },
                query: {
                    tab: 'draft',
                },
            });
        },
        handleRadar(index, row) {
            this.$router.push({
                name: 'ExperimentGroupDetail',
                params: {
                    id: row.id,
                },
                query: {
                    tab: 'radar',
                },
            });
        },
    },
};
</script>
<style lang="scss">
.el-popover__reference-wrapper {
    margin-left: 10px;
}
</style>
