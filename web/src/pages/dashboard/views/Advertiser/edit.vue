<template>
    <el-dialog :visible.sync="editFormVisible" title="编辑平台账号">
        <el-form :ref="formRef" :model="form" label-width="160px" :rules="rules">
            <el-form-item label="账号归属" label-width="160px" prop="object_id">
                <el-cascader
                    :options="objectList"
                    :props="objectProps"
                    clearable
                    placeholder="不选默认为根账号"
                    :show-all-levels="false"
                    v-model="form.object_id"
                />
            </el-form-item>
            <el-form-item label="平台类型" label-width="160px">
                <el-input v-model="form.platform_label" disabled />
            </el-form-item>
            <el-form-item label="账号名称" label-width="160px">
                <el-input v-model="form.name" disabled />
            </el-form-item>
        </el-form>
        <span slot="footer">
            <el-button @click="closeDialog">取 消</el-button>
            <el-button type="primary" @click="commit">提 交</el-button>
        </span>
    </el-dialog>
</template>
<script>
import { Advertiser, ObjectApi } from '@dashboard/api';
import { Platform } from './meta';
import Utils from '@/utils';

export default {
    name: 'AdvertiserEdit',
    inject: ['reload'],
    props: {
        defaultValue: Object,
    },
    data() {
        return {
            formRef: 'advertiserEditForm',
            editFormVisible: false,
            rules: {
                object_id: [
                    {
                        required: true,
                        message: `请选择账号归属`,
                    },
                ],
            },
            objectList: [],
            objectProps: {
                value: 'id',
                label: 'name',
                children: 'children',
                checkStrictly: true,
                emitPath: false,
            },
            form: {
                object_id: undefined,
            },
        };
    },
    watch: {
        defaultValue: {
            handler() {
                this.init();
            },
            deep: true,
        },
    },
    methods: {
        openDialog() {
            this.editFormVisible = true;
        },
        closeDialog() {
            this.editFormVisible = false;
        },
        async init() {
            await this.getObjectList();

            // this.$set(this, 'form', this.format(this.defaultValue));

            this.form = this.format(this.defaultValue);
            console.log(this.form);
        },
        getObjectList() {
            const msgPrefix = `拉取归属数据列表`;
            return new Promise((resolve, reject) => {
                if (Utils.isNonEmptyArray(this.objectList)) {
                    resolve();
                    return;
                }
                ObjectApi.get({
                    type: 'advertiser',
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
        format(data) {
            const { platform } = data || {};
            data.platform_label = Platform[platform]?.label || '未知';
            return data;
        },
        async commit() {
            console.log(this.form);
            const msgPrefix = '修改平台账号';
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    Advertiser.edit(this.form)
                        .then(() => {
                            this.$message.success(`${msgPrefix}成功`);
                            this.reload();
                        })
                        .catch((error) => {
                            console.error(error);
                            const { msg } = error || {};
                            this.$message.error(`${msgPrefix}失败 - ${msg || '未知'}`);
                        });
                }
            });
        },
    },
};
</script>
