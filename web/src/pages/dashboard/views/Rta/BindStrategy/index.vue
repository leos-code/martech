<template>
    <div>
        <el-button type="primary" @click="addDialog">
            新增策略
        </el-button>

        <bind-strategy-list @edit="editDialog" />

        <el-dialog :visible.sync="dialogFormVisible" :title="DIALOG_CONFIG.TITLE[dialogType]" @close="closeDialog">
            <bind-strategy-edit
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
import BindStrategyList from './list';
import BindStrategyEdit from './edit';
import { DIALOG_TYPE } from '@dashboard/const/dialog';

const DIALOG_CONFIG = {
    TITLE: {
        [DIALOG_TYPE.ADD]: '新增策略',
        [DIALOG_TYPE.EDIT]: '编辑策略',
    },
    COMMIT: {
        [DIALOG_TYPE.ADD]: '创 建',
        [DIALOG_TYPE.EDIT]: '编 辑',
    },
};

export default {
    name: 'BindStrategy',
    components: {
        BindStrategyList,
        BindStrategyEdit,
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
            console.log('addDialog');
            this.dialogType = DIALOG_TYPE.ADD;
            this.dialogDefaultValue = {};
            this.dialogFormVisible = true;
        },
        closeDialog() {
            console.log('closeDialog');
            this.dialogFormVisible = false;
        },
        editDialog(data) {
            console.log('editDialog');
            console.log(data);
            this.dialogType = DIALOG_TYPE.EDIT;
            this.dialogDefaultValue = data;
            this.dialogFormVisible = true;
        },
    },
};
</script>
