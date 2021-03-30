<template>
    <div>
        <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane
                v-for="objectType in objectTypes"
                :key="objectType.type"
                :name="objectType.type"
                :label="objectType.label"
            >
                <el-button type="primary" @click="addDialog(objectType)"> 新增{{ objectType.label }} </el-button>

                <object-list :object-type="objectType" @create="addDialog" @edit="editDialog" />

                <el-dialog
                    :visible.sync="objectType.dialogFormVisible"
                    :title="`${DIALOG_CONFIG.TITLE[objectType.dialogType]}${objectType.label}`"
                    @close="closeDialog(objectType)"
                    :close-on-click-modal="false"
                >
                    <object-edit
                        :object-type="objectType"
                        :type="objectType.dialogType"
                        :commit-text="DIALOG_CONFIG.COMMIT[objectType.dialogType]"
                        :default-value="objectType.dialogDefaultValue"
                        :parent-value="objectType.dialogParentValue"
                        @close="closeDialog(objectType)"
                    />
                </el-dialog>
            </el-tab-pane>
        </el-tabs>
    </div>
</template>
<script>
import _ from 'lodash';
import ObjectList from './list';
import ObjectEdit from './edit';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import Utils from '@/utils';
import { ObjectTypes } from './meta';

const DIALOG_CONFIG = {
    TITLE: {
        [DIALOG_TYPE.ADD]: '新增',
        [DIALOG_TYPE.EDIT]: '编辑',
    },
    COMMIT: {
        [DIALOG_TYPE.ADD]: '新 增',
        [DIALOG_TYPE.EDIT]: '编 辑',
    },
};

const DIALOG_DEFAULT_STATUS = {
    dialogType: DIALOG_TYPE.ADD,
    dialogFormVisible: false,
    dialogDefaultValue: {},
    dialogParentValue: {},
};

export default {
    name: 'Object',
    components: {
        ObjectList,
        ObjectEdit,
    },
    data() {
        return {
            objectTypes: [],
            DIALOG_CONFIG,
            activeName: undefined,
        };
    },
    mounted() {
        this.init();
    },
    methods: {
        init() {
            const objectTypes = ObjectTypes.map((objectType) => {
                return _.merge({}, objectType, DIALOG_DEFAULT_STATUS);
            });
            this.objectTypes = objectTypes;
            this.activeName = this.objectTypes[0]?.type;
        },
        addDialog(objectType, data) {
            objectType.dialogType = DIALOG_TYPE.ADD;
            objectType.dialogDefaultValue = {};
            // object是实体数据，则当前为 创建子角色 行为
            if (!Utils.is(data, 'MouseEvent')) {
                objectType.dialogParentValue = data;
            } else {
                objectType.dialogParentValue = {};
            }
            objectType.dialogFormVisible = true;
        },
        closeDialog(objectType) {
            objectType.dialogFormVisible = false;
            objectType.dialogDefaultValue = {};
            objectType.dialogParentValue = {};
        },
        editDialog(objectType, data) {
            objectType.dialogType = DIALOG_TYPE.EDIT;
            objectType.dialogDefaultValue = data;
            objectType.dialogParentValue = {};
            objectType.dialogFormVisible = true;
        },
        handleClick() {},
    },
};
</script>
