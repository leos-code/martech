<template>
    <div>
        <el-popover placement="top" width="160" v-model="promptPopoverVisible" v-if="canPrompt">
            <p>确定要启动吗？</p>
            <div style="text-align: right; margin: 0">
                <el-button size="mini" type="text" @click="promptPopoverVisible = false">
                    取消
                </el-button>
                <el-button type="primary" size="mini" @click="prompt">
                    确定
                </el-button>
            </div>
            <el-button type="primary" slot="reference">
                启动
            </el-button>
        </el-popover>

        <el-button @click="edit">
            编辑
        </el-button>

        <experiment-group-detail-page type="draft" :experiment-group="experimentGroup" />
    </div>
</template>
<script>
import ExperimentGroupDetailPage from '../page';
import { ExperimentGroup } from '@dashboard/api';

export default {
    name: 'ExperimentGroupDetailDraft',
    components: {
        ExperimentGroupDetailPage,
    },
    props: {
        experimentGroup: Object,
    },
    data() {
        return {
            canPrompt: false,
            promptPopoverVisible: false,
        };
    },
    watch: {
        experimentGroup: {
            handler() {
                this.refreshPrompt();
            },
            deep: true,
        },
    },
    mounted() {
        this.refreshPrompt();
    },
    methods: {
        refreshPrompt() {
            if (typeof this.experimentGroup === 'undefined') {
                return;
            }
            const { draft, current } = this.experimentGroup || {};
            const { version: draftVersion } = draft || {};
            const { version: currentVersion } = current || {};
            if (typeof currentVersion === 'undefined' || draftVersion !== currentVersion) {
                this.canPrompt = true;
            } else {
                this.canPrompt = false;
            }
        },
        prompt() {
            this.promptPopoverVisible = false;
            const { id } = this.experimentGroup || {};
            const msgPrefix = '启动实验';
            ExperimentGroup.prompt(id)
                .then(() => {
                    this.$message.success(`${msgPrefix}成功`);
                    this.$emit('refresh');
                    this.$router.push({
                        name: 'ExperimentGroupDetail',
                        params: {
                            id,
                        },
                        query: {
                            tab: 'current',
                        },
                    });
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix} - ${msg}`);
                });
        },
        edit() {
            const { id } = this.experimentGroup || {};
            this.$router.push({
                name: 'ExperimentGroupEdit',
                params: {
                    id,
                },
            });
        },
    },
};
</script>
