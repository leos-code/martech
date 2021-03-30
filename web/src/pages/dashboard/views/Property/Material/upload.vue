<template>
    <el-dialog title="素材上传详情" :visible.sync="visible" width="800px" :append-to-body="true" :modal-append-to-body="true"
               :close-on-click-modal="false" @closed="onClose">
        <div>
            <el-table :data="fileList" class="material_upload_table" height="500">
                <el-table-column label="素材名称" width="200" prop="name"></el-table-column>
                <el-table-column label="文件大小" width="100">
                    <template slot-scope="scope">
                        {{formatSize(scope.row.size)}}
                    </template>
                </el-table-column>
                <el-table-column label="上传进度">
                    <template slot-scope="scope">
                        <template  v-if="scope.row.percentage != 100">
                            <el-progress :percentage="scope.row.percentage" :show-text="false" class="material_upload_progress"></el-progress><i class="ml5 el-icon-loading material_progress_icon blue"></i>
                        </template>
                        <template v-else>
                            <template v-if="scope.row.response && scope.row.response.code == 0">
                                <el-progress :percentage="100" :show-text="false" status="success" class="material_upload_progress"></el-progress><i class="ml5 el-icon-success material_progress_icon green"></i>
                            </template>
                            <template v-else>
                                <el-progress :percentage="100" :show-text="false" status="exception" class="material_upload_progress"></el-progress>
                                <el-tooltip class="item" effect="dark" :content="formatErrmsg(scope.row)" placement="bottom">
                                    <i class="ml5 el-icon-warning material_progress_icon red"></i>
                                </el-tooltip>
                            </template>
                        </template>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </el-dialog>
</template>

<script>
    export default {
        name: "reject",
        props: {
            fileList: Array,
        },
        data(){
            return {
                rejectReason: '',
                visible: false,
                resolveFunc: '',
            }
        },
        methods: {
            open(){
                this.visible = true;
            },
            onClose(){

            },
            stringify(str){
                return JSON.stringify(str)
            },
            formatSize(size) {
                size = parseInt(size / 1024);
                if (size < 1024) {
                    return _.round(size, 2) + ' KB';
                } else {
                    size = parseInt(size) / 1024;
                    return _.round(size, 2) + ' MB';
                }
            },
            formatErrmsg(file){
                return file.response && file.response.msg || '上传失败';
            }
        }
    }
</script>

<style scoped>
    .material_upload_table{
        font-size: 13px;
    }
    .dialog-footer{
        text-align: center;
    }
    .material_upload_table .material_upload_progress{
        width: 300px;
        display: inline-block;
    }
    .material_upload_table .material_progress_icon{
        font-size: 16px;
    }
    .material_upload_table .green{
        color:#67C23A;
    }
    .material_upload_table .red{
        color: #F56C6C;
    }
    .material_upload_table .blue{
        color: #409EFF;
    }
</style>
