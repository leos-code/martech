<template>
    <div v-if="activeExperimentStage.length > 0">
        <el-tabs type="card" v-model="activeStage" @tab-click="handleClick">
            <el-tab-pane
                v-for="stage in activeExperimentStage"
                :key="stage.id"
                :label="`放量${stage.flowRate}%`"
                :name="stage.id + ''"
            >
                <experiment-group-detail-radar-query :group="experimentGroup" :stage="stage" />
                <experiment-group-detail-radar-table :group="experimentGroup" :stage="stage" />
            </el-tab-pane>
        </el-tabs>
    </div>
    <div v-else>当前无运行和历史实验，请启动后再查询效果数据</div>
</template>
<script>
import _ from 'lodash';
import ExperimentGroupDetailRadarQuery from './query';
import ExperimentGroupDetailRadarTable from './table';
import { getSpecStatusStage } from '../utils';
import { EXPERIMENT_STAGE_STATUS } from '../../meta';
import Utils from '@/utils';

export default {
    name: 'ExperimentGroupDetailRadar',
    props: {
        experimentGroup: Object,
    },
    components: {
        ExperimentGroupDetailRadarQuery,
        ExperimentGroupDetailRadarTable,
    },
    data() {
        return {
            activeExperimentStage: [],
            activeStage: '',
        };
    },
    watch: {
        experimentGroup: {
            handler() {
                this.refreshActiveStage();
            },
            deep: true,
        },
    },
    mounted() {
        this.refreshActiveStage();
    },
    methods: {
        refreshActiveStage() {
            this.activeExperimentStage = []
                .concat(getSpecStatusStage(this.experimentGroup, EXPERIMENT_STAGE_STATUS.Running.key))
                .concat(getSpecStatusStage(this.experimentGroup, EXPERIMENT_STAGE_STATUS.Stop.key));
            if (this.activeExperimentStage.length === 0) {
                return;
            }
            this.activeExperimentStage = _.orderBy(this.activeExperimentStage, ['start_time'], ['desc']);
            this.activeExperimentStage.forEach((stage) => {
                stage.flowRate = 0;
                const firstItemRtaExp = _.get(stage, ['experiment_item', 0, 'rta_exp']);
                if (Utils.isNonEmptyArray(firstItemRtaExp)) {
                    firstItemRtaExp.forEach((rtaExp) => {
                        stage.flowRate += rtaExp.flow_rate;
                    });
                }
            });
            this.activeStage = this.activeExperimentStage[this.activeExperimentStage.length - 1].id + '';
        },
        handleClick() {},
    },
};
</script>
