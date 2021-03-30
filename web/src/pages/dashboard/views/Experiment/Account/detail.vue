<template>
    <div>
        <el-card class="account-table-expand">
            <div slot="header">
                <span>基础信息</span>
            </div>
            <el-form label-position="right">
                <el-form-item label="账号ID">
                    <span>{{ accountDetail.id }}</span>
                </el-form-item>
                <el-form-item label="RTA ID">
                    <span>{{ accountDetail.rta_id }}</span>
                </el-form-item>
                <el-form-item label="Token">
                    <span>{{ accountDetail.token }}</span>
                </el-form-item>
                <el-form-item label="描述">
                    <span>{{ accountDetail.description }}</span>
                </el-form-item>
                <el-form-item label="修改时间">
                    <span>{{ accountDetail.updated_at }}</span>
                </el-form-item>
            </el-form>
        </el-card>
        <template v-for="item in accountDetail.rta_exp">
            <el-card class="account-rta-exp-item" :key="item.id">
                <div slot="header" class="clearfix">
                    <span>Rta Exp 实验组 {{ item.id }}</span>
                </div>
                <el-form label-position="right" inline>
                    <el-form-item label="Rta Exp 实验ID">
                        <span>{{ item.id }}</span>
                    </el-form-item>
                    <el-form-item label="采样比例">
                        <span>{{ item.flow_rate }}%</span>
                    </el-form-item>
                    <el-form-item label="过期时间">
                        <span>{{ item.expiration_time }}</span>
                    </el-form-item>
                    <el-form-item label="绑定状态">
                        <span :class="item.bind_status_style">{{ item.bind_status_format }}</span>
                    </el-form-item>
                    <el-form-item label="可用状态">
                        <span :class="item.status_style">{{ item.status }}</span>
                    </el-form-item>
                    <el-form-item label="同步时间">
                        <span>{{ item.updated_at }}</span>
                    </el-form-item>
                </el-form>
            </el-card>
        </template>
    </div>
</template>
<script>
import _ from 'lodash';
import { ExperimentAccount } from '@dashboard/api';
import { DateTimeFormat } from '@/utils';
import { RtaExpBindStatus, RtaExpStatus } from '@dashboard/views/Experiment/meta';

export default {
    name: 'RtaAccountDetail',
    props: {
        account: Object,
    },
    data() {
        return {
            accountDetail: {},
        };
    },
    mounted() {
        this.getRtaExpList(this.account);
    },
    methods: {
        getRtaExpList(account) {
            const msgPrefix = '拉取账户信息详情';
            if (typeof account !== 'object') {
                this.$message.error('账户信息有误, 请检查或刷新页面');
            } else {
                ExperimentAccount.getRtaExpList(account)
                    .then(({ data }) => {
                        this.accountDetail = this.format(data);
                    })
                    .catch(({ msg }) => {
                        this.$message.error(`${msgPrefix}失败 - ${msg}`);
                    });
            }
        },
        format(account) {
            if (typeof account !== 'object') {
                this.$message({
                    message: 'Rta 账户数据有误',
                    type: 'error',
                });
                return {};
            }
            const newAccount = _.clone(account);
            const { updated_at, rta_exp } = newAccount || {};
            newAccount.updated_at = DateTimeFormat.formatByDate(updated_at, 'yyyy-mm-dd HH:MM');
            if (_.isArray(rta_exp)) {
                rta_exp.map((item) => {
                    const { status, bind_status, expiration_time, updated_at } = item || {};
                    item.status = _.get(RtaExpStatus, `${status}.name`, '未知');
                    item.status_style = _.get(RtaExpStatus, `${status}.style`);
                    item.bind_status_format = _.get(RtaExpBindStatus, `${bind_status}.name`, '未知');
                    item.bind_status_style = _.get(RtaExpBindStatus, `${bind_status}.style`);
                    item.expiration_time = DateTimeFormat.formatByDate(expiration_time, 'yyyy-mm-dd HH:MM');
                    item.updated_at = DateTimeFormat.formatByDate(updated_at, 'yyyy-mm-dd HH:MM');
                });
                newAccount.rta_exp = rta_exp;
            }
            return newAccount;
        },
    },
};
</script>
<style lang="scss">
.account-table-expand {
    label {
        width: 90px;
        color: #99a9bf;
    }
    .el-form-item {
        margin-left: 10%;
        margin-right: 0;
        margin-bottom: 0;
        width: 80%;
    }
}
.account-rta-exp-item {
    display: inline-block;
    margin-right: 10px;
    width: 49%;
    label {
        width: 120px;
        color: #99a9bf;
    }
    .el-form-item {
        margin-right: 0;
        margin-bottom: 0;
        width: 80%;
    }
}
</style>
