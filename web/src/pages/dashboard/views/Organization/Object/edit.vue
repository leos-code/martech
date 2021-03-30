<template>
    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
        <el-form-item label="父角色" label-width="160px">
            <el-cascader
                :options="objectList"
                :props="parentObjectProps"
                clearable
                :placeholder="`不选默认为根${objectType.label}`"
                :show-all-levels="false"
                v-model="form.parent_id"
            />
        </el-form-item>
        <el-form-item :label="`${objectType.label}名称`" label-width="160px" prop="name">
            <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item>
            <el-button @click="close">
                取 消
            </el-button>
            <el-button type="primary" @click="commit">
                {{ commitText }}
            </el-button>
        </el-form-item>
    </el-form>
</template>
<script>
import _ from 'lodash';
import { ObjectApi } from '@dashboard/api';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import Utils from '@/utils';

const DEFAULT_FORM = {
    name: '',
    parent_id: 0,
};

export default {
    name: 'ObjectEdit',
    inject: ['reload'],
    props: {
        type: {
            type: String,
            default: DIALOG_TYPE.ADD,
        },
        commitText: String,
        defaultValue: {
            type: Object,
            default: () => {
                return {};
            },
        },
        parentValue: {
            type: Object,
            default: () => {
                return {};
            },
        },
        objectType: Object,
    },
    data() {
        return {
            formRef: 'objectEditForm',
            form: _.cloneDeep(DEFAULT_FORM),
            rules: {
                name: [
                    {
                        required: true,
                        message: `请输入${this.objectType.label}名称`,
                    },
                ],
            },
            objectList: [],
            parentObjectProps: {
                value: 'id',
                label: 'name',
                children: 'children',
                checkStrictly: true,
            },
            DIALOG_TYPE,
        };
    },
    watch: {
        defaultValue() {
            this.init();
        },
    },
    mounted() {
        this.init();
    },
    methods: {
        async init() {
            if (typeof this.defaultValue === 'undefined') {
                this.reset();
            } else {
                this.form = _.assign({}, DEFAULT_FORM, this.defaultValue);
            }

            await this.getObjectList();

            await this.getDefaultRootObject();

            // 初始化默认选中的父节点
            if (this.type === DIALOG_TYPE.ADD && Utils.isPositiveInteger(this.parentValue?.id)) {
                this.form.parent_id = this.parentValue.id;
            }
        },
        getDefaultRootObject() {
            const msgPrefix = '拉取根Object';
            return new Promise((resolve, reject) => {
                ObjectApi.get({
                    type: 'object',
                })
                    .then(({ data }) => {
                        if (Utils.isNonEmptyArray(data)) {
                            const { id } = data[0] || {};
                            this.form.object = `object_${id}`;
                        }
                        resolve();
                    })
                    .catch((error) => {
                        console.error(error);
                        const { msg } = error || {};
                        this.$message.error(`${msgPrefix}失败 - ${msg || '未知'}`);
                        reject();
                    });
            });
        },
        getObjectList() {
            const msgPrefix = `拉取${this.objectType?.label}列表`;
            return new Promise((resolve, reject) => {
                if (Utils.isNonEmptyArray(this.objectList)) {
                    resolve();
                    return;
                }
                ObjectApi.get({
                    type: this.objectType?.type,
                })
                    .then(({ data }) => {
                        this.objectList = Utils.arrayToTree(_.cloneDeep(data), {
                            parentIdNullVal: 0,
                        });
                        resolve();
                    })
                    .catch((error) => {
                        const { msg } = error || {};
                        this.$message.error(`${msgPrefix}失败 - ${msg || '未知'}`);
                        console.error(error);
                        reject();
                    });
            });
        },
        reset() {
            this.form = _.cloneDeep(DEFAULT_FORM);
            this.$refs[this.formRef].resetFields();
        },
        close() {
            this.$emit('close');
            this.reset();
        },
        async commit() {
            const msgPrefix =
                this.type === DIALOG_TYPE.ADD ? `新建${this.objectType.label}` : `修改${this.objectType.label}`;

            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    // 将parent_id列表转为单数字
                    if (Utils.isArray(this.form.parent_id)) {
                        this.form.parent_id = this.form.parent_id[0];
                    }
                    this.form.type = this.objectType.type;
                    console.log('this.form');
                    console.log(this.form);
                    ObjectApi.edit(this.form)
                        .then(() => {
                            this.$message.success(`${msgPrefix}成功`);
                            this.reload();
                        })
                        .catch((error) => {
                            const { msg } = error || {};
                            console.error(error);
                            this.$message.error(`${msgPrefix}失败 - ${msg}`);
                        });
                } else {
                    return false;
                }
            });
        },
    },
};
</script>
