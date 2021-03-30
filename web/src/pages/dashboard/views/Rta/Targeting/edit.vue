<template>
    <el-dialog
        title="RTA策略配置"
        :visible.sync="visible"
        width="1000px"
        :append-to-body="true"
        :modal-append-to-body="true"
        @opened="onOpen"
        @closed="onClose"
        v-loading="loading"
    >
        <base-info :default-value="form" ref="base_info" />
        <portrait :default-value="form.targeting_info" ref="portrait" />
        <strategy :default-value="form.bind_strategy" ref="strategy" />
        <frequency :default-value="form.freq_control.rules" ref="frequency" />
        <div slot="footer">
            <el-button @click="visible=false">
                取消
            </el-button>
            <el-button type="primary" @click="commit">
                提交
            </el-button>
        </div>
    </el-dialog>
</template>

<script>
    import BaseInfo from './BaseInfo'
    import Portrait from "./portrait";
    import Strategy from "./strategy";
    import Frequency from "./frequency";
    import _ from 'lodash';
    import {Targeting} from '@dashboard/api';

    const DEFAULT_FORM = {
        name: '',
        targeting_info: [],
        bind_strategy: [],
        freq_control: {
            rules: [],
        },
    }

    export default {
        name: "Edit",
        components: {
            BaseInfo,
            Portrait,
            Strategy,
            Frequency
        },
        props: {
            defaultValue: {
                type: Object,
                default: () => {
                    return {};
                },
            }
        },
        data() {
            return {
                visible: false,
                loading: false,
                form: _.cloneDeep(DEFAULT_FORM)
            }
        },
        methods: {
            open(){
                this.visible = true;
            },
            onOpen(){
                this.init()
            },
            onClose(){
                this.form = _.cloneDeep(DEFAULT_FORM)
            },
            init(){
                this.form = _.assign(_.cloneDeep(DEFAULT_FORM), this.defaultValue || {})
            },
            async commit(){
                let data = {}
                let baseInfoRef = this.$refs['base_info'];
                let portraitRef = this.$refs['portrait'];
                let strategyRef = this.$refs['strategy'];
                let frequencyRef = this.$refs['frequency'];
                Promise.all([baseInfoRef.validate(), portraitRef.validate(), strategyRef.validate(), frequencyRef.validate()]).then(async (pass)=>{
                    if(pass[0] && pass[1] && pass[2] && pass[3]){
                        data = Object.assign(baseInfoRef.getData(), {
                            targeting_info: portraitRef.getData(),
                            bind_strategy: strategyRef.getData(),
                            freq_control: {rules: frequencyRef.getData()},
                        });
                        let submitData = Object.assign({}, _.cloneDeep(this.defaultValue), data);
                        this.loading = true;
                        let response = await Targeting.edit(submitData);
                        this.loading = false;
                        let { data, code, msg } = response || {};
                        if(!code) {
                            this.visible = false;
                            this.$message({
                                type: 'success',
                                message: '提交成功'
                            });
                            this.$emit('success')
                        }else{
                            this.$message.error(`提交失败 - ${msg}`);
                        }
                    }else{
                        this.$message.error(`填写信息存在错误，请修正后再提交`);
                    }
                }).catch((e)=>{
                    this.$message.error(`填写信息存在错误，请修正后再提交`);
                });
            }
        }
    }
</script>

<style scoped>
    .w500 {
        width: 500px;
    }
    .w100{
        width: 100px;
    }
    .red{
        color:#F56C6C;
    }
    .gray{
        color:#999999;
    }
</style>
