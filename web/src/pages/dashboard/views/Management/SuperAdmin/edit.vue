<template>
    <el-form :ref="formRef" :model="form" label-width="120px">
        <user-search :ref="userSearchRef" label-width="120px" @select="handleUserSelected" />
        <el-form-item>
            <el-button @click="close">
                取 消
            </el-button>
            <el-button type="primary" @click="commit" :disabled="!userSelected">
                提交
            </el-button>
        </el-form-item>
    </el-form>
</template>
<script>
import _ from 'lodash';
import { UserSearch } from '@dashboard/views/Organization';
import { SuperAdmin } from '@dashboard/api';
import Utils from '@/utils';

export default {
    name: 'SuperAdminEdit',
    inject: ['reload'],
    components: {
        UserSearch,
    },
    data() {
        return {
            formRef: 'SuperAdminEditForm',
            userSearchRef: 'userSearchRef',
            form: {
                id: '',
                phone_number: '',
            },
            userSelected: false,
        };
    },
    methods: {
        close() {
            this.$emit('close');
        },
        reset() {
            this.$refs[this.userSearchRef].resetFields();
        },
        handleUserSelected(user) {
            const { id, phone_number } = user || {};
            if (Utils.isPositiveInteger(id) && Utils.isNotEmptyString(phone_number)) {
                this.form.id = id;
                this.form.phone_number = phone_number;
                this.userSelected = true;
            }
        },
        commit() {
            const msgPrefix = '新增超级管理员';
            SuperAdmin.edit(this.form)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.reset();
                    this.reload();
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
    },
};
</script>
