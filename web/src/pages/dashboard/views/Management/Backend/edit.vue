<template>
    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
        <el-form-item label="Path" label-width="160px" prop="path">
            <el-input v-model="form.path" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Method" label-width="160px" prop="method">
            <el-select v-model="form.method" placeholder="请选择Method">
                <el-option
                    v-for="apiMethod in ApiMethods"
                    :key="apiMethod.method"
                    :label="apiMethod.method"
                    :value="apiMethod.method"
                />
            </el-select>
        </el-form-item>
        <el-form-item label="组" label-width="160px" prop="group">
            <el-input v-model="form.group" autocomplete="off" />
        </el-form-item>
        <el-form-item label="描述" label-width="160px" prop="description">
            <el-input v-model="form.description" autocomplete="off" />
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
import { Backend } from '@dashboard/api';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import { ApiMethods } from './meta';

const DEFAULT_FORM = {
    method: '',
    path: '',
    group: '',
    description: '',
};

export default {
    name: 'BackendEdit',
    inject: ['reload'],
    props: {
        type: {
            type: String,
            default: DIALOG_TYPE.ADD,
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
            formRef: 'backendEditForm',
            form: _.clone(DEFAULT_FORM),
            rules: {
                method: [
                    {
                        required: true,
                        message: '请选择Method',
                        trigger: 'change',
                    },
                ],
                path: [
                    {
                        required: true,
                        message: 'Path不能为空',
                    },
                ],
                group: [
                    {
                        required: true,
                        message: '请输入所在组',
                    },
                ],
                description: [
                    {
                        required: true,
                        message: '请输入描述',
                    },
                ],
            },
            ApiMethods: Object.values(ApiMethods),
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
            const msgPrefix = this.type === DIALOG_TYPE.ADD ? '新建Api' : '修改Api';
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    Backend.edit(this.form)
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
