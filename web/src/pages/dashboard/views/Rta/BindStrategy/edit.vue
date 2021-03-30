<template>
    <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules">
        <el-form-item label="定向ID" label-width="160px" prop="targeting_id">
            <el-select size="small" v-model="form.targeting_id" placeholder="请选择定向ID">
                <template v-for="targeting in targetingList">
                    <el-option :key="targeting.id" :label="targeting.name" :value="targeting.id" />
                </template>
            </el-select>
        </el-form-item>
        <el-form-item label="策略名称" label-width="160px" prop="name">
            <el-input size="small" v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="媒体平台" label-width="160px" prop="platform">
            <el-select size="small" v-model="form.platform" placeholder="请选择媒体平台">
                <template v-for="item in PLATFORM">
                    <el-option :key="item.key" :label="item.name" :value="item.key" />
                </template>
            </el-select>
        </el-form-item>
        <el-form-item label="RTA策略ID (可选)" label-width="160px" prop="strategy_id">
            <ug-auto-split-input size="small" :value="form.strategy.strategy_id" @input="inputStrategyId" />
        </el-form-item>
        <el-form-item label="广告主ID (可选)" label-width="160px" prop="advertiser_id">
            <ug-auto-split-input size="small" :value="form.strategy.advertiser_id" @input="inputAdvertiserId" />
        </el-form-item>
        <el-form-item label="推广计划ID (可选)" label-width="160px" prop="campaign_id">
            <ug-auto-split-input size="small" :value="form.strategy.campaign_id" @input="inputCampaignId" />
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
import UgAutoSplitInput from '@/components/UgAutoSplitInput';
import { BindStrategy, Targeting } from '@dashboard/api';
import { DIALOG_TYPE } from '@dashboard/const/dialog';
import PLATFORM from '@dashboard/const/platform';

const DEFAULT_FORM = {
    name: '',
    platform: '',
    strategy: {
        strategy_id: [],
        advertiser_id: [],
        campaign_id: [],
    },
};

export default {
    name: 'BindStrategyEdit',
    inject: ['reload'],
    components: {
        UgAutoSplitInput,
    },
    props: {
        type: {
            type: String,
            default: () => {
                return DIALOG_TYPE.ADD;
            },
        },
        title: {
            type: String,
            default: () => {
                return '';
            },
        },
        commitText: {
            type: String,
            default: () => {
                return '';
            },
        },
        defaultValue: {
            type: Object,
            default: () => {
                return {};
            },
        },
    },
    watch: {
        defaultValue() {
            this.init();
        },
    },
    data() {
        return {
            formRef: 'bindStrategyEditForm',
            PLATFORM,
            targetingList: [],
            form: _.clone(DEFAULT_FORM),
            rules: {
                targeting_id: [
                    {
                        required: true,
                        message: '请选择定向ID',
                        trigger: 'change',
                    },
                ],
                name: [
                    {
                        required: true,
                        message: '请输入策略名称',
                        trigger: 'blur',
                    },
                    {
                        min: 3,
                        max: 30,
                        message: '名称长度在 3 到 30 个字符',
                        trigger: 'blur',
                    },
                ],
                platform: [
                    {
                        required: true,
                        message: '请选择平台',
                        trigger: 'change',
                    },
                ],
                strategy_id: [],
                advertiser_id: [],
                campaign_id: [],
            },
        };
    },
    mounted() {
        this.init();
    },
    async created() {
        await this.initTargeting();
    },
    methods: {
        init() {
            if (typeof this.defaultValue === 'undefined') {
                console.log('empty defaultValue');
            }
            this.form = _.assign({}, DEFAULT_FORM, this.defaultValue);
        },
        async initTargeting() {
            const response = await Targeting.getList();
            const { status, data: responseData } = response || {};
            const { code, data } = responseData || {};
            if (status === 200 && code === 0) {
                this.targetingList = data;
            }
        },
        close() {
            this.$emit('close');
        },
        async commit() {
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    const response = await BindStrategy.edit(this.form);
                    const { data } = response || {};
                    const { code, msg } = data || {};
                    if (code === 0) {
                        this.$message({
                            message: '新建成功',
                            type: 'success',
                        });
                        this.reload();
                    } else {
                        this.$message.error(msg);
                    }
                } else {
                    return false;
                }
            });
        },
        inputStrategyId(strategyId) {
            this.form.strategy.strategy_id = strategyId;
        },
        inputAdvertiserId(advertiserId) {
            this.form.strategy.advertiser_id = advertiserId;
        },
        inputCampaignId(campaignId) {
            this.form.strategy.campaign_id = campaignId;
        },
    },
};
</script>
