<template>
    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
        <el-form-item label="租户名称" label-width="160px" prop="name">
            <el-input v-model="form.name" autocomplete="off" @keyup.enter.native="commit" />
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
import { Tenant } from '@dashboard/api';
import { DIALOG_TYPE } from '@dashboard/const/dialog';

const DEFAULT_FORM = {
    name: '',
};

export default {
    name: 'TenantEdit',
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
            formRef: 'tenantEditForm',
            form: _.clone(DEFAULT_FORM),
            rules: {
                name: [
                    {
                        required: true,
                        message: '请输入租户名称',
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
            const msgPrefix = this.type === DIALOG_TYPE.ADD ? '新建租户' : '修改租户';
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    Tenant.edit(this.form)
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
