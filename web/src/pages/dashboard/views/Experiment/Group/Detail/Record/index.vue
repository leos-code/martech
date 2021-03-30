<template>
    <div class="experiment-modify-records">
        <el-row :gutter="24">
            <el-col :span="18" :offset="1">
                <el-timeline>
                    <el-timeline-item
                        v-for="(record, index) in modifyRecords"
                        :key="index"
                        size="large"
                        :timestamp="record.updated_at_format"
                        :icon="record.operation.icon"
                        :color="record.operation.color"
                        placement="top"
                    >
                        <template v-if="record.operation.key === OPERATION.Modify.key">
                            <el-collapse>
                                <el-collapse-item
                                    :title="
                                        `${record.operation.name}（ ${record.user ? record.user.name : 'Unknown'} ）`
                                    "
                                >
                                    <template
                                        v-if="
                                            record.diff &&
                                                record.diff.visualizeResult &&
                                                record.diff.visualizeResult.length > 0
                                        "
                                    >
                                        <el-tree
                                            :data="record.diff.visualizeResult"
                                            :props="diffTreeProps"
                                            default-expand-all
                                            :render-content="renderTree"
                                        />
                                    </template>
                                </el-collapse-item>
                            </el-collapse>
                        </template>
                        <template v-else>
                            <span class="operation-name">
                                {{ record.operation.name }}
                                （{{ record.user ? record.user.name : 'Unknown' }}）
                            </span>
                        </template>
                    </el-timeline-item>
                </el-timeline>
            </el-col>
        </el-row>
    </div>
</template>
<script>
import _ from 'lodash';
import Utils from '@/utils';
import { OPERATION } from './meta';
import Diff from './diff';

export default {
    name: 'ExperimentGroupDetailRecord',
    props: {
        experimentGroup: {
            type: Object,
            default: () => {
                return {};
            }
        },
    },
    data() {
        return {
            modifyRecords: [],
            OPERATION,
            diffTreeProps: {
                label: 'label',
                children: 'children',
            },
        };
    },
    watch: {
        experimentGroup: {
            handler() {
                this.filterRecords();
            },
            deep: true,
        },
    },
    mounted() {
        this.filterRecords();
    },
    methods: {
        filterRecords() {
            const { modify_record = [] } = this.experimentGroup || {};
            let originModifyRecords = _.cloneDeep(modify_record);
            // 将修改记录按照ID倒序排序，用于列表展示
            originModifyRecords = _.orderBy(originModifyRecords, ['id'], ['desc']);
            // 获取Modify记录的对比数据，并处理对比结果可视化
            originModifyRecords.forEach((record, index) => {
                if (record.operation === OPERATION.Modify.key) {
                    record.diff = {
                        base: [],
                        item: [],
                    };
                    const lhsRecord = _.find(
                        originModifyRecords,
                        {
                            operation: OPERATION.Modify.key,
                        },
                        index + 1,
                    );
                    if (lhsRecord) {
                        record.diff = Diff.diffRecords(_.cloneDeep(lhsRecord), _.cloneDeep(record));
                    }
                }
            });

            // 将第一个Modify状态置为Create
            const firstModifyIndex = _.findLastIndex(originModifyRecords, {
                operation: OPERATION.Modify.key,
            });
            if (firstModifyIndex !== -1) {
                originModifyRecords[firstModifyIndex].operation = OPERATION.Create.key;
            }
            // 单个record记录的字段可视化
            originModifyRecords.forEach((record) => {
                record.updated_at_format = Utils.DateTimeFormat.formatByDate(record.updated_at, 'yyyy-mm-dd HH:MM');
                // 将operation转成对象，方便视图层处理和展示
                record.operation = OPERATION[record.operation];
            });

            this.modifyRecords = originModifyRecords;
        },
        renderTree(h, { data }) {
            const { label, patchType } = data || {};
            const { class: patchClass, label: patchLabel } = patchType || {};
            if (patchClass) {
                return (
                    <span>
                        <span class={patchClass}>{patchLabel}：</span>
                        {label}
                    </span>
                );
            } else {
                return <span>{label}</span>;
            }
        },
    },
};
</script>
<style lang="scss">
.experiment-modify-records {
    margin-top: 20px;
    .operation-name {
        font-size: 13px;
        font-weight: 500;
    }
    .operation-user {
        font-size: 13px;
    }
}
</style>
