<template>
    <div>
        <el-button type="primary" @click="addDialog">
            新增功能
        </el-button>

        <el-button type="success" @click="openSyncDialog">
            发布功能
        </el-button>

        <feature-list @edit="editDialog" />

        <el-dialog :visible.sync="dialogFormVisible" :title="DIALOG_CONFIG.TITLE[dialogType]" @close="closeDialog">
            <feature-edit
                :type="dialogType"
                :title="DIALOG_CONFIG.TITLE[dialogType]"
                :commit-text="DIALOG_CONFIG.COMMIT[dialogType]"
                :default-value="dialogDefaultValue"
                @close="closeDialog"
            />
        </el-dialog>

        <el-dialog :visible.sync="dialogSyncVisible" title="确认发布功能">
            <span>请确认是否要发布功能配置？发布后将针对所有租户生效</span>
            <span slot="footer">
                <el-button @click="closeSyncDialog">取 消</el-button>
                <el-button type="primary" @click="sync">确 定</el-button>
            </span>
        </el-dialog>
    </div>
</template>
<script>
import FeatureList from './list';
import FeatureEdit from './edit';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import { Feature } from '@dashboard/api';

const DIALOG_CONFIG = {
    TITLE: {
        [DIALOG_TYPE.ADD]: '新增功能',
        [DIALOG_TYPE.EDIT]: '编辑功能',
    },
    COMMIT: {
        [DIALOG_TYPE.ADD]: '创 建',
        [DIALOG_TYPE.EDIT]: '提 交',
    },
};

export default {
    name: 'Feature',
    components: {
        FeatureList,
        FeatureEdit,
    },
    data() {
        return {
            DIALOG_CONFIG,
            dialogType: DIALOG_TYPE.ADD,
            dialogFormVisible: false,
            dialogDefaultValue: {},
            dialogSyncVisible: false,
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
        openSyncDialog() {
            this.dialogSyncVisible = true;
        },
        closeSyncDialog() {
            this.dialogSyncVisible = false;
        },
        sync() {
            const msgPrefix = '发布功能';
            Feature.sync()
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.dialogSyncVisible = false;
                })
                .catch((error) => {
                    console.error(error);
                    const { msg } = error || {};
                    this.$message.success(`${msgPrefix}失败 - ${msg || '未知'}`);
                    this.dialogSyncVisible = false;
                });
        },
    },
};
</script>
