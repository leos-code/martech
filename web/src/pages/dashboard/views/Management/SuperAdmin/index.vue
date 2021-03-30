<template>
    <div>
        <el-button type="primary" @click="addDialog">
            新增超级管理员
        </el-button>

        <super-admin-list />

        <el-dialog :visible.sync="dialogFormVisible" :title="DIALOG_CONFIG.TITLE[dialogType]" @close="closeDialog">
            <super-admin-edit
                :ref="editRef"
                :type="dialogType"
                :title="DIALOG_CONFIG.TITLE[dialogType]"
                :commit-text="DIALOG_CONFIG.COMMIT[dialogType]"
                :default-value="dialogDefaultValue"
                @close="closeDialog"
            />
        </el-dialog>
    </div>
</template>
<script>
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import SuperAdminList from './list';
import SuperAdminEdit from './edit';

const DIALOG_CONFIG = {
    TITLE: {
        [DIALOG_TYPE.ADD]: '新增 超级管理员',
    },
    COMMIT: {
        [DIALOG_TYPE.ADD]: '创 建',
    },
};

export default {
    name: 'Admin',
    components: {
        SuperAdminList,
        SuperAdminEdit,
    },
    data() {
        return {
            DIALOG_CONFIG,
            dialogType: DIALOG_TYPE.ADD,
            dialogFormVisible: false,
            dialogDefaultValue: {},
            editRef: 'editSuperAdminRef',
        };
    },
    methods: {
        addDialog() {
            this.dialogType = DIALOG_TYPE.ADD;
            this.dialogDefaultValue = {};
            this.dialogFormVisible = true;
        },
        closeDialog() {
            this.dialogFormVisible = false;
            this.$refs[this.editRef].reset();
            this.dialogDefaultValue = {};
        },
    },
};
</script>
