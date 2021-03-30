<template>
    <el-card class="box-card mt50">
        <div slot="header" class="clearfix">
            <span>用户画像</span>
        </div>
        <div>
            <el-tabs v-model="activeName" type="border-card">
                <el-tab-pane label="包含" name="include">
                    <targeting ref="includeRef" :schema="schema" :default-value="includeValues" :not="false" @active="active" />
                </el-tab-pane>
                <el-tab-pane label="排除" name="exclude">
                    <targeting ref="excludeRef" :schema="schema" :default-value="excludeValues" :not="true" @active="active" />
                </el-tab-pane>
            </el-tabs>
        </div>
    </el-card>
</template>

<script>

    import {Schema} from '@dashboard/api';
    import Targeting from "./targeting";
    import _ from 'lodash';


    const INTEGER_VALUE = {start: '', end: ''};

    export default {
        name: "Portrait",
        components: {
            Targeting
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
                activeName: 'include',
                includeValues: [],
                excludeValues: [],
                schema: [],
            }
        },
        methods: {
            async getSchema() {
                const response = await Schema.getSchema();
                const {code, data, msg} = response || {};
                if (!code) {
                    this.schema = data && data.fields || [];
                } else {
                    this.schema = [];
                    this.$message.error(msg)
                }
            },
            //由targeting组件触发外部选中该tab,点击提交触发校验以后，自动跳转到校验不通过的tab来
            active(not){
                if(not){
                    this.activeName = 'exclude';
                }else{
                    this.activeName = 'include';
                }
            },
            validate() {
                return new Promise((resolve, reject) => {
                    let includeRef = this.$refs['includeRef'];
                    let excludeRef = this.$refs['excludeRef'];
                    Promise.all([includeRef.validate(), excludeRef.validate()]).then((pass) => {
                        resolve(pass.findIndex(item => !item) === -1);
                    }).catch((e) => {
                        reject(e);
                    });
                });
            },
            getData() {
                let includeRef = this.$refs['includeRef'];
                let excludeRef = this.$refs['excludeRef'];
                return _.concat(includeRef.getData(), excludeRef.getData());
            },

            // 区分包含和排除
            classify(){
                this.includeValues = [];
                this.excludeValues = [];
                (this.defaultValue || []).forEach((item)=>{
                    if(item.not){
                        this.excludeValues.push(item);
                    }else{
                        this.includeValues.push(item);
                    }
                })
            },
            init() {
                this.activeName = 'include';
                this.classify();
            },
        },
        async mounted() {
            this.init();
            await this.getSchema();
        }
    }
</script>

<style scoped>
    .w500 {
        width: 500px;
    }

    .w100 {
        width: 100px;
    }

    .red {
        color: #F56C6C;
    }

    .gray {
        color: #999999;
    }
</style>
