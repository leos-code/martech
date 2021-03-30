<template>
    <div v-if="canShow">
        <el-popover v-model="stopPopoverVisible" placement="top" width="160">
            <p>确定要废弃吗？</p>
            <div style="text-align: right; margin: 0">
                <el-button size="mini" type="text" @click="stopPopoverVisible = false">
                    取消
                </el-button>
                <el-button type="primary" size="mini" @click="stop">
                    确定
                </el-button>
            </div>
            <el-button type="primary" slot="reference">
                停止
            </el-button>
        </el-popover>

        <experiment-group-detail-page type="current" :experiment-group="experimentGroup" />
    </div>
    <div v-else>
        当前无线上版本，请前往 实验管理 启动实验
    </div>
</template>
<script>
import ExperimentGroupDetailPage from '../page';
import { ExperimentGroup } from '@dashboard/api';

export default {
    name: 'ExperimentGroupDetailCurrent',
    components: {
        ExperimentGroupDetailPage,
    },
    props: {
        experimentGroup: Object,
    },
    data() {
        return {
            canShow: false,
            stopPopoverVisible: false,
        };
    },
    watch: {
        experimentGroup: {
            handler() {
                this.initShowStatus();
            },
            deep: true,
        },
    },
    mounted() {
        this.initShowStatus();
    },
    methods: {
        initShowStatus() {
            const { current } = this.experimentGroup || {};
            if (typeof current !== 'undefined') {
                this.canShow = true;
            } else {
                this.canShow = false;
            }
        },
        stop() {
            this.stopPopoverVisible = false;
            const { id } = this.experimentGroup || {};
            const msgPrefix = '停止实验';
            ExperimentGroup.stop(id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.$emit('refresh');
                    this.$router.push({
                        name: 'ExperimentGroupDetail',
                        params: {
                            id,
                        },
                        query: {
                            tab: 'draft',
                        },
                    });
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix} - ${msg}`);
                });
        },
    },
};
</script>
