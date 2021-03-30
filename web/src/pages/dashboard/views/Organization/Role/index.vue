<template>
    <div>
        <el-button type="primary" @click="addDialog">
            新增角色
        </el-button>

        <role-list @create="addDialog" @edit="editDialog" />

        <el-dialog
            :visible.sync="dialogFormVisible"
            :title="DIALOG_CONFIG.TITLE[dialogType]"
            @close="closeDialog"
            :close-on-click-modal="false"
        >
            <role-edit
                :type="dialogType"
                :commit-text="DIALOG_CONFIG.COMMIT[dialogType]"
                :default-value="dialogDefaultValue"
                :parent-value="dialogParentValue"
                @close="closeDialog"
            />
        </el-dialog>
    </div>
</template>
<script>
import RoleList from './list';
import RoleEdit from './edit';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import Utils from '@/utils';

const DIALOG_CONFIG = {
    TITLE: {
        [DIALOG_TYPE.ADD]: '新增角色',
        [DIALOG_TYPE.EDIT]: '编辑角色',
    },
    COMMIT: {
        [DIALOG_TYPE.ADD]: '新 增',
        [DIALOG_TYPE.EDIT]: '编 辑',
    },
};

export default {
    name: 'Role',
    components: {
        RoleList,
        RoleEdit,
    },
    data() {
        return {
            DIALOG_CONFIG,
            dialogType: DIALOG_TYPE.ADD,
            dialogFormVisible: false,
            dialogDefaultValue: {},
            dialogParentValue: {},
        };
    },
    methods: {
        addDialog(data) {
            this.dialogType = DIALOG_TYPE.ADD;
            this.dialogDefaultValue = {};
            // role是实体数据，则当前为 创建子角色 行为
            if (!Utils.is(data, 'MouseEvent')) {
                this.dialogParentValue = data;
            } else {
                this.dialogParentValue = {};
            }
            this.dialogFormVisible = true;
        },
        closeDialog() {
            this.dialogFormVisible = false;
            this.dialogDefaultValue = {};
            this.dialogParentValue = {};
        },
        editDialog(data) {
            this.dialogType = DIALOG_TYPE.EDIT;
            this.dialogDefaultValue = data;
            this.dialogParentValue = {};
            this.dialogFormVisible = true;
        },
    },
};
</script>
