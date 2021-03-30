<template>
    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
        <el-form-item label="RTA ID" label-width="160px" prop="rta_id">
            <el-input autocomplete="off" v-model="form.rta_id" />
        </el-form-item>
        <el-form-item label="Token" label-width="160px" prop="token">
            <el-input autocomplete="off" v-model="form.token" />
        </el-form-item>
        <el-form-item label="描述" label-width="160px" prop="description">
            <el-input autocomplete="off" v-model="form.description" @keyup.enter.native="commit" />
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
import { ExperimentAccount } from '@dashboard/api';
import { DIALOG_TYPE } from '@dashboard/const/dialog';

const DEFAULT_FORM = {
    rta_id: '',
    token: '',
    description: '',
};

export default {
    name: 'ExperimentAccountEdit',
    inject: ['reload'],
    props: {
        type: {
            type: String,
            default: () => {
                return DIALOG_TYPE.ADD;
            },
        },
        title: String,
        commitText: String,
        defaultValue: Object,
    },
    watch: {
        defaultValue() {
            this.init();
        },
    },
    data() {
        return {
            formRef: 'experimentAccountEditForm',
            form: _.clone(DEFAULT_FORM),
            rules: {
                rta_id: [
                    {
                        required: true,
                        message: '请输入RTA ID',
                        trigger: 'change',
                    },
                ],
                token: [
                    {
                        required: true,
                        message: '请输入Token',
                        trigger: 'change',
                    },
                ],
            },
        };
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
            const msgPrefix = this.type === DIALOG_TYPE.ADD ? '创建账户' : '修改账户';
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    ExperimentAccount.edit(this.form)
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
