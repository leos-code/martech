<template>
    <div>
        <div class="ml10 mr10">
            <el-input v-model="query.name" placeholder="输入名称查询" size="small" class="w200 mr15" prefix-icon="el-icon-search" @change="getTargetingList(1)"></el-input>

            <el-button type="primary" class="fr" @click="handleEdit(0, {})" size="small"><i class="el-icon-circle-plus"></i>新增策略</el-button>
        </div>

        <el-divider></el-divider>

        <div v-loading="loading">

            <el-table ref="strategyTable" stripe :data="list">
                <el-table-column prop="id" label="策略Id" sortable column-key="id" width="200" />
                <el-table-column prop="name" label="策略名称" />
                <el-table-column prop="created_at" label="创建日期" :formatter="dateFormat" />
                <el-table-column prop="updated_at" label="最近更新日期" :formatter="dateFormat" />
                <el-table-column label="操作" width="200px">
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
            <el-pagination
                class="pagination"
                @current-change="handleCurrentChange"
                :current-page="0"
                :page-size="pageParams.page_size"
                layout="total, prev, pager, next"
                :total="total">
            </el-pagination>

        </div>


        <edit ref="edit" :default-value="editDefaultValue" @success="getTargetingList"></edit>
    </div>
</template>
<script>
    import Edit from './edit';
    import {Targeting} from '@dashboard/api';

    export default {
        name: 'Targeting',
        components: {
            Edit,
        },
        data(){
            return {
                loading: false,
                query: {
                    name: '',
                },
                total: 0,
                pageParams: {
                    page: 1,
                    page_size: 20,
                },
                list: [],
                editDefaultValue:{}
            }
        },
        methods: {
            getQueryField(){
                let queryFields = [];
                if(this.query.name){
                    queryFields.push({
                        field: 'name',
                        value: this.query.name,
                        operation: 'like'
                    })
                }
                return queryFields
            },
            async getTargetingList(page){
                const msgPrefix = '拉取策略配置列表';
                if(page !== undefined){
                    this.pageParams.page = page;
                }
                this.loading = true;
                let filter = this.getQueryField();
                const response = await Targeting.getList(Object.assign({}, {filter: encodeURIComponent(JSON.stringify(filter))}, this.pageParams));
                const { data, code, msg } = response || {};
                this.loading = false;
                if (!code) {
                    this.total = data.total || 0;
                    this.list = data.list || [];
                } else {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                }
            },

            dateFormat(row, column, cellValue){
                return new Date(cellValue).toLocaleString()
            },

            handleEdit(idx, value){
                this.editDefaultValue = value;
                this.$refs['edit'].open();
            },
            async handleDelete(index, row) {
                row.visible = false;
                const response = await Targeting.delete(row.id);
                const { code, msg } = response || {};
                if (!code) {
                    this.$message({
                        message: '删除成功',
                        type: 'success',
                    });
                    this.getTargetingList(this.pageParams.page);
                } else {
                    this.$message.error(`删除失败 - ${msg}`);
                }
            },
            //换页操作
            handleCurrentChange(v) {
                this.getTargetingList(v);
            },
        },
        mounted(){
            this.getTargetingList();
        }
    };
</script>
<style scoped>
    .fr{
        float: right
    }
</style>
