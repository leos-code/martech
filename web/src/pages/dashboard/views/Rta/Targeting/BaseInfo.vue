<template>
    <el-card class="box-card">
        <div slot="header" class="clearfix">
            <span>基础信息</span>
        </div>
        <div>
            <el-form ref="formRef" :model="form" label-width="100px" :rules="rules">
                <el-form-item prop="name">
                    <span slot="label" class="gray mr10">名称：</span>
                    <el-input class="w500" size="small" v-model="form.name" autocomplete="off" />
                </el-form-item>
            </el-form>
        </div>
    </el-card>
</template>

<script>

    import _ from "lodash";

    const DEFAULT_FORM = {
        name: '',
    };
    export default {
        name: "BaseInfo",
        props: {
            defaultValue: {
                type: Object,
                default: () => {
                    return {};
                },
            },
        },
        watch: {
            defaultValue() {
                this.init();
            }
        },
        data() {
            return {
                form: {},
                rules: {
                    name: [
                        { required: true, message: '请输入RTA策略名称', trigger: 'blur' }
                    ],
                }
            }
        },
        methods: {
            validate(){
                return new Promise((resolve, reject)=>{
                    this.$refs['formRef'].validate().then((pass)=>{
                        resolve(pass)
                    }).catch((e)=>{
                        reject(e)
                    });
                });
            },
            getData(){
                return _.cloneDeep(this.form);
            },
            init(){
                this.form = _.assign({}, _.cloneDeep(DEFAULT_FORM), _.cloneDeep(this.defaultValue))
                this.$refs['formRef'].clearValidate();
            },
        },
        mounted(){
            this.init();
        }
    }
</script>

<style scoped>
    .gray{
        color:#999;
    }
</style>
