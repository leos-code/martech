<template>
    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
        <el-form-item label="名称" label-width="160px" prop="name">
            <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="描述" label-width="160px">
            <el-input v-model="form.description" autocomplete="off" />
        </el-form-item>
        <el-form-item label="配置" label-width="160px">
            <el-tabs type="border-card">
                <el-tab-pane label="页面菜单">
                    <el-tree
                        :ref="menuSelectRef"
                        :data="menuList"
                        show-checkbox
                        node-key="id"
                        default-expand-all
                        :props="menuTreeProps"
                        check-on-click-node
                        :expand-on-click-node="false"
                        :check-strictly="true"
                    />
                </el-tab-pane>
                <el-tab-pane label="页面子功能">
                    <el-tree
                        :ref="subFunctionSelectRef"
                        :data="subFunctionList"
                        show-checkbox
                        node-key="id"
                        default-expand-all
                        :props="subFunctionTreeProps"
                        check-on-click-node
                        :expand-on-click-node="false"
                        :check-strictly="true"
                    />
                </el-tab-pane>
                <el-tab-pane label="Api">
                    <el-tree
                        :ref="apiSelectRef"
                        :data="apiList"
                        show-checkbox
                        node-key="id"
                        default-expand-all
                        :props="apiTreeProps"
                        check-on-click-node
                        :expand-on-click-node="false"
                        :check-strictly="true"
                    />
                </el-tab-pane>
            </el-tabs>
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
import { Feature, Frontend, Backend } from '@dashboard/api';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import { FrontendTypes } from '../Frontend/meta';
import Utils from '@/utils';

const DEFAULT_FORM = {
    name: '',
    description: '',
};

export default {
    name: 'FeatureEdit',
    inject: ['reload'],
    props: {
        type: {
            type: String,
            default: DIALOG_TYPE.ADD,
        },
        commitText: {
            type: String,
            default: '',
        },
        defaultValue: {
            type: Object,
            default: () => {
                return {};
            },
        },
    },
    data() {
        return {
            formRef: 'featureEditForm',
            form: _.clone(DEFAULT_FORM),
            rules: {
                name: [
                    {
                        required: true,
                        message: '请输入功能名称',
                    },
                ],
            },
            frontendList: [],
            menuList: [],
            menuSelectRef: 'menuSelectRef',
            menuTreeProps: {
                children: 'children',
                label: 'description',
            },
            subFunctionList: [],
            subFunctionSelectRef: 'subFunctionSelectRef',
            subFunctionTreeProps: {
                children: 'children',
                label: 'description',
            },
            apiList: [],
            apiSelectRef: 'apiSelectRef',
            apiTreeProps: {
                children: 'children',
                label: 'description',
            },
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
            console.log('init called');
            if (typeof this.defaultValue === 'undefined') {
                this.reset();
            } else {
                this.form = _.assign({}, DEFAULT_FORM, this.defaultValue);
            }

            await this.getFrontendList();
            this.setDefaultSelectedFrontends();
            await this.getBackendList();
            this.setDefaultSelectedBackends();
        },
        getFrontendList() {
            const msgPrefix = `拉取页面菜单列表`;
            return new Promise((resolve, reject) => {
                if (Utils.isNonEmptyArray(this.frontendList)) {
                    resolve();
                    return;
                }
                Frontend.get()
                    .then(({ data }) => {
                        this.frontendList = _.cloneDeep(data);
                        this.menuList = _.filter(this.frontendList, {
                            type: FrontendTypes.Menu.type,
                        });
                        this.subFunctionList = _.filter(this.frontendList, {
                            type: FrontendTypes.SubFunction.type,
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
        getSelectedFrontends() {
            const selectedMenus = this.$refs[this.menuSelectRef].getCheckedNodes();
            const selectedSubFunctions = this.$refs[this.subFunctionSelectRef].getCheckedNodes();
            return selectedMenus.concat(selectedSubFunctions);
        },
        setDefaultSelectedFrontends() {
            const { frontend: frontends } = this.form || {};
            let defaultSelectedMenus = [];
            let defaultSelectedSubFunctions = [];
            if (Utils.isNonEmptyArray(frontends)) {
                defaultSelectedMenus = _.filter(frontends, {
                    type: FrontendTypes.Menu.type,
                });
                defaultSelectedSubFunctions = _.filter(frontends, {
                    type: FrontendTypes.SubFunction.type,
                });
            }
            this.$refs[this.menuSelectRef].setCheckedNodes(defaultSelectedMenus);
            this.$refs[this.subFunctionSelectRef].setCheckedNodes(defaultSelectedSubFunctions);
        },
        getBackendList() {
            const msgPrefix = `拉取Api列表`;
            return new Promise((resolve, reject) => {
                if (Utils.isNonEmptyArray(this.apiList)) {
                    resolve();
                    return;
                }
                Backend.get()
                    .then(({ data }) => {
                        this.apiList = _.cloneDeep(data);
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
        getSelectedBackends() {
            return this.$refs[this.apiSelectRef].getCheckedNodes();
        },
        setDefaultSelectedBackends() {
            const { backend = [] } = this.form || {};
            this.$refs[this.apiSelectRef].setCheckedNodes(backend);
        },
        reset() {
            this.$refs[this.formRef].resetFields();
        },
        close() {
            this.$emit('close');
        },
        async commit() {
            const msgPrefix = this.type === DIALOG_TYPE.ADD ? '新建功能' : '修改功能';
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    Feature.edit(this.form)
                        .then(({ data }) => {
                            data.frontend = this.getSelectedFrontends();
                            data.backend = this.getSelectedBackends();
                            Feature.relate(data).then(() => {
                                this.$message.success(`${msgPrefix}成功`);
                                this.reload();
                            });
                        })
                        .catch((error) => {
                            const { msg } = error || {};
                            console.error(error);
                            this.$message.error(`${msgPrefix}失败 - ${msg || '未知'}`);
                        });
                } else {
                    return false;
                }
            });
        },
    },
};
</script>
