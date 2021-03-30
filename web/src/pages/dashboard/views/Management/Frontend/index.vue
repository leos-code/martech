<template>
    <div>
        <el-button type="primary" @click="addDialog">
            新增菜单
        </el-button>

        <frontend-list @edit="editDialog" />

        <el-dialog :visible.sync="dialogFormVisible" :title="DIALOG_CONFIG.TITLE[dialogType]" @close="closeDialog">
            <frontend-edit
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
import FrontendList from './list';
import FrontendEdit from './edit';
import { DIALOG_TYPE } from '@dashboard/const/dialog';

const DIALOG_CONFIG = {
    TITLE: {
        [DIALOG_TYPE.ADD]: '新增菜单/子功能',
        [DIALOG_TYPE.EDIT]: '编辑菜单/子功能',
    },
    COMMIT: {
        [DIALOG_TYPE.ADD]: '创 建',
        [DIALOG_TYPE.EDIT]: '提 交',
    },
};

export default {
    name: 'Frontend',
    components: {
        FrontendList,
        FrontendEdit,
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
        },
        editDialog(data) {
            this.dialogType = DIALOG_TYPE.EDIT;
            this.dialogDefaultValue = data;
            this.dialogFormVisible = true;
        },
    },
};
</script>
