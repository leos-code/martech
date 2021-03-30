<template>
    <el-row :gutter="24">
        <el-col :span="18" :offset="2">
            <el-form label-position="right">
                <el-divider content-position="left">
                    基本信息
                </el-divider>
                <el-form-item label="组织名称" label-width="160px">
                    {{ info.name }}
                </el-form-item>
                <el-form-item label="组织ID" label-width="160px">
                    {{ info.id }}
                </el-form-item>
            </el-form>
        </el-col>
    </el-row>
</template>
<script>
import { OrganizationInfo } from '@dashboard/api';

export default {
    name: 'OrganizationInfo',
    data() {
        return {
            info: {},
        };
    },
    mounted() {
        this.init();
    },
    methods: {
        init() {
            const msgPrefix = '拉取组织信息';
            OrganizationInfo.get()
                .then(({ data }) => {
                    this.info = data;
                })
                .catch((error) => {
                    console.error(error);
                    const { msg } = error || {};
                    this.$message.error(`${msgPrefix}失败 - ${msg || '未知'}`);
                });
        },
    },
};
</script>
