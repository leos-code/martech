<template>
    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
        <el-form-item label="参数名称" label-width="160px" prop="name">
            <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="参数类型" label-width="160px" prop="type">
            <el-select v-model="form.type" placeholder="请选择参数类型">
                <template v-for="item in PARAMETER_TYPE">
                    <el-option :key="item.key" :label="item.name" :value="item.key" />
                </template>
            </el-select>
        </el-form-item>
        <el-form-item label="描述" label-width="160px" prop="description">
            <el-input v-model="form.description" autocomplete="off" @keyup.enter.native="commit" />
        </el-form-item>
        <el-form-item>
            <el-button @click="close">
                取 消
            </el-button>
            <el-button type="primary" @click="commit">
                {{ commitText }}
            </el-button>
        </el-form-item>
    </el-form>
</template>
<script>
import _ from 'lodash';
import { ExperimentParameter } from '@dashboard/api';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import { PARAMETER_TYPE } from './meta';

const DEFAULT_FORM = {
    name: '',
    type: '',
    description: '',
};

export default {
    name: 'ExperimentParameterEdit',
    inject: ['reload'],
    props: {
        type: {
            type: String,
            default: DIALOG_TYPE.ADD,
        },
        title: {
            type: String,
            default: '',
        },
        commitText: {
            type: String,
            default: '',
        },
        defaultValue: {
            type: Object,
            default: () => {
                return {};
            },
        },
    },
    data() {
        return {
            formRef: 'experimentParameterEditForm',
            PARAMETER_TYPE,
            form: _.clone(DEFAULT_FORM),
            rules: {
                name: [
                    {
                        required: true,
                        message: '请输入参数名称',
                        trigger: 'change',
                    },
                ],
                type: [
                    {
                        required: true,
                        message: '请选择参数类型',
                        trigger: 'change',
                    },
                ],
            },
        };
    },
    watch: {
        defaultValue() {
            this.init();
        },
    },
    mounted() {
        this.init();
    },
    methods: {
        init() {
            if (typeof this.defaultValue === 'undefined') {
                this.reset();
            } else {
                this.form = _.assign({}, DEFAULT_FORM, this.defaultValue);
            }
        },
        reset() {
            this.$refs[this.formRef].resetFields();
        },
        close() {
            this.$emit('close');
        },
        async commit() {
            const msgPrefix = this.type === DIALOG_TYPE.ADD ? '新建参数' : '修改参数';
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    ExperimentParameter.edit(this.form)
                        .then(() => {
                            this.$message.success(`${msgPrefix}成功`);
                            this.reload();
                        })
                        .catch(({ msg }) => {
                            this.$message.error(`${msgPrefix}失败 - ${msg}`);
                        });
                } else {
                    return false;
                }
            });
        },
    },
};
</script>
