<template>
    <el-tabs v-model="activeName" @tab-click="handleClick">
        <el-tab-pane label="效果报表" name="radar">
            <experiment-group-detail-radar :experiment-group="experimentGroup" />
        </el-tab-pane>
        <el-tab-pane label="线上版本" name="current">
            <experiment-group-detail-current :experiment-group="experimentGroup" @refresh="refresh" />
        </el-tab-pane>
        <el-tab-pane label="实验管理" name="draft">
            <experiment-group-detail-draft :experiment-group="experimentGroup" @refresh="refresh" />
        </el-tab-pane>
        <el-tab-pane label="修改历史" name="record">
            <experiment-group-detail-record :experiment-group="experimentGroup" />
        </el-tab-pane>
    </el-tabs>
</template>
<script>
import ExperimentGroupDetailDraft from './Draft';
import ExperimentGroupDetailCurrent from './Current';
import ExperimentGroupDetailRecord from './Record';
import ExperimentGroupDetailRadar from './Radar';
import { ExperimentGroup } from '@dashboard/api';
import { isNotEmptyString } from '@/utils';

export default {
    name: 'ExperimentGroupDetail',
    components: {
        ExperimentGroupDetailDraft,
        ExperimentGroupDetailCurrent,
        ExperimentGroupDetailRecord,
        ExperimentGroupDetailRadar,
    },
    data() {
        return {
            activeName: 'radar',
            experimentGroup: undefined,
        };
    },
    watch: {
        $route() {
            this.getActiveTab();
        },
    },
    mounted() {
        this.getActiveTab();
        this.getGroup();
    },
    methods: {
        handleClick(tab) {
            const { name } = tab || {};
            if (name) {
                this.$router.push({
                    path: this.$router.history.current.path,
                    query: {
                        tab: name,
                    },
                });
            }
        },
        getActiveTab() {
            const { tab } = this.$router.history.current.query || {};
            if (isNotEmptyString(tab)) {
                this.activeName = tab;
            }
        },
        getGroup() {
            const id = this.$route.params.id;
            const msgPrefix = '查询实验信息';
            ExperimentGroup.get(id)
                .then(({ data }) => {
                    this.experimentGroup = data;
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix} - ${msg}`);
                    this.experimentGroup = undefined;
                });
        },
        refresh() {
            this.getGroup();
        },
    },
};
</script>
