<template>
    <div>
        <el-button type="primary" @click="addDialog">
            新增租户
        </el-button>

        <tenant-list @edit="editDialog" />

        <el-dialog :visible.sync="dialogFormVisible" :title="DIALOG_CONFIG.TITLE[dialogType]" @close="closeDialog">
            <tenant-edit
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
import TenantList from './list';
import TenantEdit from './edit';
import { DIALOG_TYPE } from '@dashboard/const/dialog';

const DIALOG_CONFIG = {
    TITLE: {
        [DIALOG_TYPE.ADD]: '新增租户',
        [DIALOG_TYPE.EDIT]: '编辑租户',
    },
    COMMIT: {
        [DIALOG_TYPE.ADD]: '创 建',
        [DIALOG_TYPE.EDIT]: '提 交',
    },
};

export default {
    name: 'Tenant',
    components: {
        TenantList,
        TenantEdit,
    },
    data() {
        return {
            DIALOG_CONFIG,
            dialogType: DIALOG_TYPE.ADD,
            dialogFormVisible: false,
            dialogDefaultValue: {},
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
            this.dialogDefaultValue = {};
        },
        editDialog(data) {
            this.dialogType = DIALOG_TYPE.EDIT;
            this.dialogDefaultValue = data;
            this.dialogFormVisible = true;
        },
    },
};
</script>
