<template>
    <el-card class="box-card mt50">
        <div slot="header" class="clearfix">
            <span>策略</span>
        </div>
        <div>
            <div class="mb20">
                <el-button type="primary" size="small" @click="add()">
                    新增策略
                </el-button>
            </div>
            <el-card class="strategy_card" v-for="(item, idx) in bindStrategys" :key="idx">
                <div>
                    <el-form :ref="`formRef`" :model="item" label-width="130px" :rules="rules">
                        <el-form-item label="策略名称" prop="name">
                            <el-input class="w200" size="small" v-model="item.name" autocomplete="off" />
                        </el-form-item>
                        <el-form-item label="媒体平台" prop="platform">
                            <el-select size="small" v-model="item.platform" placeholder="请选择媒体平台">
                                <template v-for="opt in PLATFORM">
                                    <el-option :key="opt.key" :label="opt.name" :value="opt.key" />
                                </template>
                            </el-select>
                        </el-form-item>
                        <el-form-item label="RTA策略ID (可选)" prop="strategy_id">
                            <ug-auto-split-input size="small" v-model="item.strategy.strategy_id" />
                        </el-form-item>
                        <el-form-item label="广告主ID (可选)" prop="advertiser_id">
                            <ug-auto-split-input size="small" v-model="item.strategy.advertiser_id" />
                        </el-form-item>
                        <el-form-item label="推广计划ID (可选)" prop="campaign_id">
                            <ug-auto-split-input size="small" v-model="item.strategy.campaign_id" />
                        </el-form-item>
                    </el-form>
                    <div class="strategy_btn_wrap">
                        <el-button class="strategy_del_btn" size="small" type="danger" @click="del(idx)">
                            删除
                        </el-button>
                    </div>
                </div>
            </el-card>
        </div>
    </el-card>
</template>

<script>
    import _ from "lodash"

    import UgAutoSplitInput from '@/components/UgAutoSplitInput';
    import PLATFORM from '@dashboard/const/platform';

    const DEFAULT_FORM = {
        name: '',
        platform: '',
        strategy: {
            strategy_id: [],
            advertiser_id: [],
            campaign_id: [],
        },
    };

    export default {
        name: "Strategy",
        components: {
            UgAutoSplitInput
        },
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
        data() {
            return {
                PLATFORM,
                bindStrategys: [],
                rules: {
                    platform: [
                        { required: true, message: '媒体平台不能为空', trigger: 'change' },
                    ]
                }
            }
        },
        methods: {
            add(){
                this.bindStrategys.push(_.cloneDeep(DEFAULT_FORM))
            },
            del(idx){
                this.bindStrategys.splice(idx, 1);
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
                return _.cloneDeep(this.bindStrategys);
            },
            init(){
                if(_.isArray(this.defaultValue)){
                    this.bindStrategys = _.cloneDeep(this.defaultValue).map((item)=>{
                        return _.merge(_.cloneDeep(DEFAULT_FORM), item)
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
    .strategy_card{
        width: 400px;
        display: inline-block;
        margin-right: 20px;
        margin-bottom: 10px;
    }
    .strategy_card .strategy_btn_wrap{
        text-align: center;
    }
</style>
