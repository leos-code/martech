<template>
    <el-dialog title="系统提示" :visible.sync="visible" width="600px" :append-to-body="true" :modal-append-to-body="true"
               :close-on-click-modal="false" @closed="onClose">
        <div class="mb10">请输入驳回原因：</div>
        <div>
            <el-input
                type="textarea"
                :rows="4"
                v-model.trim="rejectReason">
            </el-input>
        </div>
        <span slot="footer" class="dialog-footer">
            <el-button @click="visible = false">取 消</el-button>
            <el-button type="primary" @click="confirm">确认驳回</el-button>
        </span>
    </el-dialog>
</template>

<script>
    export default {
        name: "reject",
        data(){
            return {
                rejectReason: '',
                visible: false,
                resolveFunc: '',
            }
        },
        methods: {
            open(){
                return new Promise((resolve)=>{
                    this.visible = true;
                    this.resolveFunc = resolve;
                });
            },
            onClose(){
                this.rejectReason = '';
            },
            confirm(){
                if(!this.rejectReason){
                    this.$message.warning("请输入驳回原因");
                    return;
                }
                this.visible = false;
                this.resolveFunc(this.rejectReason);
            }
        }
    }
</script>

<style scoped>

</style>
