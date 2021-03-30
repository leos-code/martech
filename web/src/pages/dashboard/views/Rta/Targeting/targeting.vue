<template>
    <el-form ref="formRef" :model="schemaValue" :rules="rules">
        <el-form-item v-for="item in schema" :key="item.name" label-width="150px" :prop="item.name">
            <span slot="label" class="gray mr10">{{ item.display_name+':' }}</span>
            <template v-if="item.type === 'enum'">
                <checkbox-tree :ref="`checkbox-tree-${item.name}`" :data="item.enum" v-model="schemaValue[item.name]" />
            </template>
            <template v-if="item.type === 'integer'">
                <div v-for="(v, vidx) in schemaValue[item.name]" :key="vidx" class="mb10">
                    <el-input size="small" class="w200" v-model.number="schemaValue[item.name][vidx].start" />
                    -
                    <el-input size="small" class="w200" v-model.number="schemaValue[item.name][vidx].end" />
                    <el-button type="text" class="ml5" @click="addSchemaValue(schemaValue[item.name], vidx, item.type)">
                        <i class="el-icon-circle-plus" />
                    </el-button>
                    <el-button type="text" class="ml5 red" @click="delSchemaValue(schemaValue[item.name], vidx)" v-if="schemaValue[item.name].length > 1">
                        <i class="el-icon-error" />
                    </el-button>
                </div>
                <div class="ml5 gray">
                    (默认左闭右开)
                </div>
            </template>
            <template v-if="item.type === 'string'">
                <div v-for="(v, vidx) in schemaValue[item.name]" :key="vidx" class="mb10">
                    <el-input size="small" class="w500" v-model="schemaValue[item.name][vidx]" />
                    <el-button type="text" class="ml5" @click="addSchemaValue(schemaValue[item.name], vidx, item.type)">
                        <i class="el-icon-circle-plus" />
                    </el-button>
                    <el-button type="text" class="ml5 red" @click="delSchemaValue(schemaValue[item.name], vidx)" v-if="schemaValue[item.name].length > 1">
                        <i class="el-icon-error" />
                    </el-button>
                </div>
            </template>
        </el-form-item>
    </el-form>
</template>

<script>

    import {Schema} from '@dashboard/api';
    import CheckboxTree from "./CheckboxTree";
    import _ from 'lodash';


    const INTEGER_VALUE = {start: '', end: ''};

    export default {
        name: "Targeting",
        components: {
            CheckboxTree,
        },
        props: {
            schema: {
                type: Array,
                default: () => {
                    return [];
                },
            },
            defaultValue: {
                type: Array,
                default: () => {
                    return [];
                },
            },
            not: {
                type: Boolean,
                default: () => {
                    return false;
                },
            }
        },
        watch: {
            schema() {
                this.init();
            },
            defaultValue() {
                this.init();
            }
        },

        data() {
            return {
                schemaValue: {},
                rules: {},
            }
        },
        methods: {
            getSchemaDefaultValue(type){
                let v = this.getSchemaDefaultValueItem(type)
                return (v !== undefined ? [v] : [])
            },
            getSchemaDefaultValueItem(type) {
                if (type === 'enum') {
                    return undefined;
                } else if (type === 'string') {
                    return '';
                } else if (type === 'integer') {
                    return _.cloneDeep(INTEGER_VALUE);
                }
            },
            addSchemaValue(list, idx, type) {
                let defaultVal = this.getSchemaDefaultValueItem(type);
                if (defaultVal !== undefined) {
                    list.push(defaultVal);
                }
            },
            delSchemaValue(list, idx) {
                list.splice(idx, 1);
            },
            validate() {
                return new Promise((resolve, reject) => {
                    this.$refs['formRef'].validate().then((pass) => {
                        resolve(pass);
                    }).catch((e) => {
                        reject(e);
                        this.$emit('active', this.not)
                    });
                });
            },
            getData() {
                let value = [];
                (this.schema || []).forEach((item, idx) => {
                    let v = this.schemaValue[item.name];
                    let newObj = {
                        name: item.name,
                        not: this.not,
                        values: {}
                    };
                    if (item.type === 'enum' || item.type === 'string') {
                        v = (v || []).filter(item => !!item.trim());
                        if (v && v.length) {
                            newObj.values = {
                                type: 'string',
                                string: v
                            };
                            value.push(newObj)
                        }
                    } else {
                        v = (v || []).filter(item => item.start !== '' && item.end !== '');
                        if (v && v.length) {
                            newObj.values = {
                                type: 'range',
                                range: v
                            };
                            value.push(newObj)
                        }
                    }
                });
                return value;
            },

            rangeRules(rule, value, callback) {
                if (!_.isArray(value)) {
                    callback()
                } else {
                    let errExist = value.findIndex((item) => {
                        return !((item.start && item.end && _.isNumber(item.start) && _.isNumber(item.end)) || (!item.start && !item.end));
                    }) > -1;
                    if (errExist) {
                        callback(new Error('请填写完整和正确的数据'));
                    } else {
                        callback();
                    }
                }
            },

            init() {
                let schemaValue = {};
                let rules = {};
                (this.schema || []).forEach((item, idx) => {
                    let obj = (this.defaultValue || []).find(v => v.name === item.name);
                    if (!!obj && obj.values) {
                        let val = obj.values.type == 'range' ? obj.values.range : obj.values.string;
                        schemaValue[item.name] = val ? val : this.getSchemaDefaultValue(item.type);
                    } else {
                        schemaValue[item.name] = this.getSchemaDefaultValue(item.type);
                    }
                    //初始化rules;
                    if (item.type === 'integer') {
                        rules[item.name] = [
                            {validator: this.rangeRules, trigger: []}
                        ]
                    }
                });
                this.rules = rules;
                this.schemaValue = schemaValue;
            },
        },
        async mounted() {
            this.init();
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
