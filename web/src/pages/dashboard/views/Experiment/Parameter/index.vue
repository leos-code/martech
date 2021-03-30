<template>
    <div>
        <el-button type="primary" @click="addDialog">
            新增参数
        </el-button>

        <experiment-parameter-list @edit="editDialog" />

        <el-dialog :visible.sync="dialogFormVisible" :title="DIALOG_CONFIG.TITLE[dialogType]" @close="closeDialog">
            <experiment-parameter-edit
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
import ExperimentParameterList from './list';
import ExperimentParameterEdit from './edit';
import { DIALOG_TYPE } from '@dashboard/const/dialog';

const DIALOG_CONFIG = {
    TITLE: {
        [DIALOG_TYPE.ADD]: '新增参数',
        [DIALOG_TYPE.EDIT]: '编辑参数',
    },
    COMMIT: {
        [DIALOG_TYPE.ADD]: '创 建',
        [DIALOG_TYPE.EDIT]: '提 交',
    },
};

export default {
    name: 'ExperimentParamter',
    components: {
        ExperimentParameterList,
        ExperimentParameterEdit,
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
