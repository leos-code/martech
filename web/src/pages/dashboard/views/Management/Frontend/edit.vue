<template>
    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
        <el-form-item label="类型" label-width="160px" prop="type">
            <el-select v-model="form.type" placeholder="请选择类型">
                <el-option
                    v-for="frontendType in FrontendTypes"
                    :key="frontendType.type"
                    :label="frontendType.label"
                    :value="frontendType.type"
                />
            </el-select>
        </el-form-item>
        <el-form-item label="Key" label-width="160px" prop="key">
            <el-input v-model="form.key" autocomplete="off" />
        </el-form-item>
        <el-form-item label="组" label-width="160px" prop="group">
            <el-input v-model="form.group" autocomplete="off" />
        </el-form-item>
        <el-form-item label="描述" label-width="160px">
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
import { Frontend } from '@dashboard/api';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import { FrontendTypes } from './meta';

const DEFAULT_FORM = {
    type: undefined,
    key: '',
    group: '',
    description: '',
};

export default {
    name: 'FrontendEdit',
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
            formRef: 'frontendEditForm',
            form: _.clone(DEFAULT_FORM),
            rules: {
                type: [
                    {
                        required: true,
                        message: '请选择类型',
                        trigger: 'change',
                    },
                ],
                key: [
                    {
                        required: true,
                        message: 'Key不能为空',
                    },
                ],
                group: [
                    {
                        required: true,
                        message: '请输入所在组',
                    },
                ],
            },
            FrontendTypes: Object.values(FrontendTypes),
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
            const msgPrefix = this.type === DIALOG_TYPE.ADD ? '新建菜单/子功能' : '修改菜单/子功能';
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    Frontend.edit(this.form)
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
