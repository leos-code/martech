<template>
    <el-card class="box-card mt50">
        <div slot="header" class="clearfix">
            <span>频率控制</span>
        </div>
        <div>
            <div class="mb20">
                <el-button type="primary" size="small" @click="add()">
                    新增频控
                </el-button>
            </div>
            <el-card class="freq_card" v-for="(item, idx) in freqRules" :key="idx">
                <div>
                    <el-form ref="formRef" :model="item" label-width="100px" :rules="rules">
                        <el-form-item label="平台" prop="platform">
                            <el-select size="small" class="w250" v-model="item.platform" placeholder="请选择平台">
                                <el-option key="" value="" label="全平台" />
                                <el-option v-for="opt in PLATFORM" :key="opt.key" :value="opt.key" :label="opt.name" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="频控行为" prop="for">
                            <el-select size="small" class="w250" v-model="item.for" placeholder="请选择频控行为">
                                <el-option v-for="opt in FREQ_FOR" :key="opt.key" :value="opt.key" :label="opt.name" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="天级限制次数" prop="limit">
                            <el-input size="small" class="w250" v-model.number="item.limit" />
                        </el-form-item>
                        <el-form-item label="总共限制次数" prop="accumulate_limit">
                            <el-input size="small" class="w250" v-model.number="item.accumulate_limit" />
                        </el-form-item>
                        <el-form-item label="间隔时间" prop="interval">
                            <el-input size="small" class="w250" v-model.number="item.interval" />
                        </el-form-item>
                    </el-form>
                    <div class="freq_btn_wrap">
                        <el-button class="freq_del_btn" size="small" type="danger" @click="del(idx)">
                            删除
                        </el-button>
                    </div>
                </div>
            </el-card>
        </div>
    </el-card>
</template>

<script>
    import PLATFORM from '@dashboard/const/platform';
    import {FREQ_FOR} from '@dashboard/const/targeting';
    import _ from "lodash";

    const DEFAULT_FORM = {
        platform: '',
        for: '',
        limit: '',
        accumulate_limit: '',
        interval:''
    };

    const validateNumOrEmpty = (rule, value, callback) => {
        if(value === '' || _.isNumber(value)){
            callback()
        }else{
            callback(new Error("填写内容必须为数字"))
        }
    }

    export default {
        name: "Frequency",
        props: {
            defaultValue: {
                type: Array,
                default: () => {
                    return [];
                },
            },
        },
        watch: {
            defaultValue() {
                this.init();
            }
        },
        data(){
            return{
                PLATFORM,
                FREQ_FOR,
                freqRules:[],
                rules: {
                    for: [
                        { required: true, message: '频控行为不能为空', trigger: 'change' },
                    ],
                    limit: [
                        { validator: validateNumOrEmpty, trigger: 'blur' }
                    ],
                    accumulate_limit: [
                        { validator: validateNumOrEmpty, trigger: 'blur' }
                    ],
                    interval: [
                        { validator: validateNumOrEmpty, trigger: 'blur' }
                    ],
                }
            }
        },
        methods: {
            add(){
                this.freqRules.push(_.cloneDeep(DEFAULT_FORM))
            },
            del(idx){
                this.freqRules.splice(idx, 1);
            },
            validate(){
                return new Promise((resolve, reject)=>{
                    var validators = [];
                    var refs = this.$refs['formRef'];
                    if(!refs || !refs.length){
                        resolve(true);
                        return;
                    }
                    refs.forEach((ref)=>{
                        validators.push(ref.validate());
                    });
                    Promise.all(validators).then((pass)=>{
                        resolve(pass.findIndex(item => !item) === -1);
                    }).catch((e)=>{
                        reject(e)
                    })
                });
            },
            getData(){
                return _.cloneDeep(this.freqRules);
            },
            init(){
                if(_.isArray(this.defaultValue)){
                    this.freqRules = _.cloneDeep(this.defaultValue).map((item)=>{
                        return _.assign(_.cloneDeep(DEFAULT_FORM), item)
                    })
                }
            },
        },
        mounted() {
            this.init();
        }
    }
</script>

<style scoped>
    .freq_card{
        width: 400px;
        display: inline-block;
        margin-right: 20px;
        margin-bottom: 10px;
    }
    .freq_card .freq_btn_wrap{
        text-align: center;
    }
</style>
