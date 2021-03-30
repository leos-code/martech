<template>
    <el-row :gutter="24">
        <el-col :span="18" :offset="2">
            <el-form label-position="right">
                <el-divider content-position="left">
                    基本信息
                </el-divider>
                <el-form-item label="实验组名称" label-width="160px">
                    {{ innerExpGroup.name }}
                </el-form-item>
                <el-form-item label="实验目标" label-width="160px">
                    {{ innerExpGroup.description }}
                </el-form-item>
                <el-form-item label="RTA账户" label-width="160px">
                    {{ innerExpGroup.rta_account && innerExpGroup.rta_account.rta_id }}
                </el-form-item>
                <el-form-item label="关注人" label-width="160px">
                    <template v-for="user in innerExpGroup.user">
                        <el-tag :key="user.id">
                            {{ user.name }}
                        </el-tag>
                    </template>
                </el-form-item>
                <el-form-item label="修改时间" label-width="160px">
                    {{ innerExpStage.updated_at_format }}
                </el-form-item>
                <el-form-item label="结束时间" label-width="160px">
                    {{ innerExpStage.end_time_format }}
                </el-form-item>
                <el-divider content-position="left">
                    实验信息
                </el-divider>
                <el-tabs v-if="innerExpStage.experiment_item && innerExpStage.experiment_item.length > 0" type="card">
                    <el-tab-pane
                        v-for="exp_item in innerExpStage.experiment_item"
                        :key="exp_item.id"
                        :label="exp_item.name || `实验-${exp_item.id}`"
                    >
                        <el-form-item label="实验名称" label-width="160px">
                            {{ exp_item.name }}
                        </el-form-item>
                        <el-form-item label="RtaExpID" label-width="160px">
                            <template v-for="rta_exp_item in exp_item.rta_exp">
                                <el-tag :key="rta_exp_item.id">
                                    {{ rta_exp_item.id }}
                                </el-tag>
                            </template>
                        </el-form-item>
                        <el-form-item label="采样比例" label-width="160px">
                            {{ exp_item.total_flow_rate }}%
                        </el-form-item>
                        <el-form-item label="实验参数" label-width="160px">
                            <div v-for="metadata in exp_item.experiment_metadata" :key="metadata.id">
                                {{ metadata.experiment_parameter.name }}={{ metadata.value }}
                            </div>
                        </el-form-item>
                    </el-tab-pane>
                </el-tabs>
            </el-form>
        </el-col>
    </el-row>
</template>
<script>
import _ from 'lodash';
import Utils from '@/utils';

export default {
    name: 'ExperimentGroupDetailPage',
    props: {
        experimentGroup: {
            type: Object,
            default: () => {
                return {};
            },
        },
        type: {
            type: String,
            default: '',
        },
    },
    data() {
        return {
            innerExpGroup: {},
            innerExpStage: {},
        };
    },
    watch: {
        experimentGroup: {
            handler() {
                this.formatInnerData();
            },
            deep: true,
        },
        type() {
            this.formatInnerData();
        },
    },
    mounted() {
        this.formatInnerData();
    },
    methods: {
        formatInnerData() {
            let group = _.cloneDeep(this.experimentGroup);
            let stage = _.get(group, this.type);
            if (typeof stage === 'undefined' || typeof group === 'undefined') {
                return;
            }
            stage.updated_at_format = Utils.DateTimeFormat.formatByDate(stage.updated_at, 'yyyy-mm-dd HH:MM');
            stage.end_time_format = Utils.DateTimeFormat.formatByDate(stage.end_time, 'yyyy-mm-dd HH:MM');
            if (Utils.isNonEmptyArray(stage.experiment_item)) {
                stage.experiment_item.map((exp_item) => {
                    const { rta_exp } = exp_item || {};
                    exp_item.total_flow_rate = 0;
                    if (Utils.isNonEmptyArray(rta_exp)) {
                        rta_exp.map((rta_exp_item) => {
                            exp_item.total_flow_rate += rta_exp_item.flow_rate;
                        });
                    }
                });
            }

            this.innerExpGroup = group;
            this.innerExpStage = stage;
        },
    },
};
</script>
