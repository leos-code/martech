<template>
    <div>
        <el-button type="primary" @click="createDialog">
            新增平台账号
        </el-button>
        <el-dialog :visible.sync="createFormVisible" title="新增平台账号">
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
                <el-form-item label="平台类型" label-width="160px" prop="url">
                    <el-select v-model="form.url" placeholder="请选择平台类型">
                        <el-option
                            v-for="authority in authorityList"
                            :key="authority.platform"
                            :label="authority.label"
                            :value="authority.url"
                        />
                    </el-select>
                </el-form-item>
            </el-form>
            <span slot="footer">
                <el-button @click="closeDialog">取 消</el-button>
                <el-button type="primary" @click="commit">创 建</el-button>
            </span>
        </el-dialog>
        <el-dialog :visible.sync="authVisible" title="授权账号结果" :close-on-click-modal="false" @close="closeAuthDialog">
            若您已完成账号授权，请点击刷新列表页......
            <span slot="footer">
                <el-button @click="closeAuthDialog">刷 新</el-button>
            </span>
        </el-dialog>
    </div>
</template>
<script>
import _ from 'lodash';
import { Advertiser, ObjectApi } from '@dashboard/api';
import { Platform } from './meta';
import Utils from '@/utils';

export default {
    name: 'AdvertiserCreate',
    inject: ['reload'],
    data() {
        return {
            formRef: 'advertiserCreateForm',
            createFormVisible: false,
            form: {
                object_id: '',
                url: '',
            },
            rules: {
                object_id: [
                    {
                        required: true,
                        message: '请选择账号归属',
                    },
                ],
                url: [
                    {
                        required: true,
                        message: `请选择平台类型`,
                        trigger: 'change',
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
            authVisible: false,
            authorityList: [],
        };
    },
    methods: {
        createDialog() {
            this.createFormVisible = true;
            this.init();
        },
        closeDialog() {
            this.createFormVisible = false;
            this.reset();
        },
        reset() {
            this.$refs[this.formRef].resetFields();
        },
        async init() {
            await this.getObjectList();
            await this.getAuthorityList();
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
        getAuthorityList() {
            const msgPrefix = `拉取平台授权列表`;
            return new Promise((resolve, reject) => {
                if (Utils.isNonEmptyArray(this.authorityList)) {
                    resolve();
                    return;
                }
                Advertiser.authorize()
                    .then(({ data }) => {
                        this.authorityList = this.formatAuthority(_.cloneDeep(data));
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
        formatAuthority(data) {
            data.map((authority) => {
                const { platform } = authority || {};
                authority.label = Platform[platform]?.label || '未知';
            });
            return data;
        },
        async commit() {
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    const authUrl = this.form.url.replace(/__OBJECT_ID__/g, this.form.object_id);
                    this.authVisible = true;
                    this.closeDialog();
                    window.open(authUrl);
                } else {
                    return false;
                }
            });
        },
        closeAuthDialog() {
            this.authVisible = false;
            this.reload();
        },
    },
};
</script>
