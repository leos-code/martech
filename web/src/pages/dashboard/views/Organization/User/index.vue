<template>
    <div>
        <el-button type="primary" @click="addDialog">
            新增用户
        </el-button>

        <user-list @edit="editDialog" />

        <el-dialog
            :visible.sync="dialogFormVisible"
            :title="DIALOG_CONFIG.TITLE[dialogType]"
            @close="closeDialog"
            :close-on-click-modal="false"
        >
            <user-edit
                :type="dialogType"
                :commit-text="DIALOG_CONFIG.COMMIT[dialogType]"
                :default-value="dialogDefaultValue"
                @close="closeDialog"
            />
        </el-dialog>
    </div>
</template>
<script>
import UserList from './list';
import UserEdit from './edit';
import { DIALOG_TYPE } from '@dashboard/const/dialog';

const DIALOG_CONFIG = {
    TITLE: {
        [DIALOG_TYPE.ADD]: '绑定用户',
        [DIALOG_TYPE.EDIT]: '编辑用户',
    },
    COMMIT: {
        [DIALOG_TYPE.ADD]: '绑 定',
        [DIALOG_TYPE.EDIT]: '绑 定',
    },
};

export default {
    name: 'User',
    components: {
        UserList,
        UserEdit,
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
