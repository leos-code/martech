<template>
    <el-row :gutter="24">
        <el-col :span="18" :offset="2">
            <el-form :inline="true" :ref="formRef" :model="form" :rules="rules">
                <el-row>
                    <el-col>
                        <el-form-item label="时间粒度" label-width="100px">
                            <el-radio-group v-model="form.granularity" @change="changeGranularity">
                                <el-radio v-for="item in Granularity" :key="item.key" :label="item.key">
                                    {{ item.name }}
                                </el-radio>
                            </el-radio-group>
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-row>
                    <el-col>
                        <el-form-item label="时间范围" label-width="100px" prop="datetime_range">
                            <el-date-picker
                                v-model="form.datetime_range"
                                type="datetimerange"
                                align="right"
                                range-separator="至"
                                start-placeholder="开始日期"
                                end-placeholder="结束日期"
                                value-format="yyyyMMddHHmmss"
                                :picker-options="pickerOptions"
                            />
                        </el-form-item>
                        <el-form-item>
                            <el-checkbox
                                v-model="form.time_merge_type"
                                class="date-range-expand-checkbox"
                                :true-label="1"
                                :false-label="0"
                            >
                                累积
                            </el-checkbox>
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-row>
                    <el-col>
                        <el-form-item label="基准实验" label-width="100px" prop="base_experiment_item_id">
                            <el-select
                                v-model="form.base_experiment_item_id"
                                placeholder="请选择基准实验"
                                @change="changeBaseExp"
                            >
                                <el-option
                                    v-for="exp_item in baseExpItemArray"
                                    :key="exp_item.id"
                                    :label="exp_item.name"
                                    :value="exp_item.id"
                                />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="对照实验" label-width="100px">
                            <el-select multiple v-model="form.lab_experiment_item_id" placeholder="请选择对照实验">
                                <el-option
                                    v-for="exp_item in labExpItemArray"
                                    :key="exp_item.id"
                                    :label="exp_item.name"
                                    :value="exp_item.id"
                                    :disabled="exp_item.disable"
                                />
                            </el-select>
                        </el-form-item>
                        <el-form-item label=" " label-width="100px">
                            <el-button type="primary" @click="handleQuery">
                                查询
                            </el-button>
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-row>
                    <el-col>
                        <el-collapse>
                            <el-collapse-item title="更多查询条件">
                                <el-form-item label="广告主ID" label-width="100px">
                                    <ug-auto-split-input v-model="form.filter.uid" />
                                </el-form-item>
                                <el-form-item label="APP ID" label-width="100px">
                                    <ug-auto-split-input v-model="form.filter.appid" />
                                </el-form-item>
                                <el-form-item label="推广计划ID" label-width="100px">
                                    <ug-auto-split-input v-model="form.filter.cid" />
                                </el-form-item>
                                <el-form-item label="广告ID" label-width="100px">
                                    <ug-auto-split-input v-model="form.filter.aid" />
                                </el-form-item>
                                <el-form-item label="用户权重" label-width="100px">
                                    <ug-auto-split-input v-model="form.filter.user_weight" />
                                </el-form-item>
                            </el-collapse-item>
                        </el-collapse>
                    </el-col>
                </el-row>
            </el-form>
        </el-col>
    </el-row>
</template>
<script>
import _ from 'lodash';
import UgAutoSplitInput from '@/components/UgAutoSplitInput';
import { Granularity } from './meta';

/**
 * TODO list
 *  - feature: 时间选择器的粒度变化时，支持已选时间联动
 *  - feature: 时间选择器支持按照当前阶段的起始时间，动态修改可选时间窗口
 *  - feature: 确定接口后，修改时间选择器的datetime_range映射值
 */
export default {
    name: 'ExperimentGroupDetailRadarQuery',
    props: {
        group: Object,
        stage: Object,
    },
    components: {
        UgAutoSplitInput,
    },
    computed: {
        formRef() {
            return `radar-form-group${this.group.id}-stage-${this.stage.id}`;
        },
        radarId() {
            return `radar-group${this.group.id}-stage${this.stage.id}`;
        },
        pickerOptions() {
            const granularity = _.find(this.Granularity, {
                label: this.form.granularity,
            });
            const max = new Date().getTime();
            let min = -Infinity;
            if (typeof granularity !== 'undefined' && _.isSafeInteger(granularity.range)) {
                min = max - 86400000 * granularity.range;
            }
            let disabledDate = (time) => {
                return time.getTime() < min || time.getTime() > max;
            };
            return {
                disabledDate,
            };
        },
        baseExpItemArray() {
            return _.cloneDeep(this.stage.experiment_item);
        },
        labExpItemArray() {
            return _.cloneDeep(this.stage.experiment_item);
        },
    },
    data() {
        return {
            form: {
                granularity: 'day',
                time_merge_type: 0,
                lab_experiment_item_id: [],
                filter: {},
                datetime_range: [],
            },
            rules: {
                datetime_range: [
                    {
                        type: 'array',
                        required: true,
                        message: '请选择时间区间',
                        trigger: 'change',
                    },
                ],
                base_experiment_item_id: [
                    {
                        required: true,
                        message: '请选择基准实验',
                        trigger: 'change',
                    },
                ],
            },
            datePickerOptions: {},
            Granularity,
        };
    },
    beforeDestroy() {
        this.$store.commit('experiment/radar/delStageMap', {
            radarId: this.radarId,
        });
    },
    beforeUpdate() {
        this.initRadarStore(this.radarId);
    },
    beforeMount() {
        this.initRadarStore(this.radarId);
    },
    methods: {
        handleQuery() {
            this.$refs[this.formRef].validate((valid) => {
                if (valid) {
                    const queryOptions = _.cloneDeep(this.form);
                    const { datetime_range = [] } = queryOptions || {};
                    queryOptions.begin_time = _.toNumber(datetime_range[0]);
                    queryOptions.end_time = _.toNumber(datetime_range[1]);
                    queryOptions.experiment_stage_id = this.stage.id;
                    this.$store.commit('experiment/radar/activeQuerying', {
                        radarId: this.radarId,
                        options: queryOptions,
                    });
                } else {
                    return false;
                }
            });
        },
        changeGranularity() {
            this.$set(this.form, 'datetime_range', []);
        },
        changeBaseExp(baseExpItemId) {
            // 重置disable状态
            this.labExpItemArray.map((expItem) => delete expItem.disable);
            // 将已选中的基准实验组置灰
            const disabledExpItemIndex = _.findIndex(this.labExpItemArray, {
                id: baseExpItemId,
            });
            this.labExpItemArray[disabledExpItemIndex].disable = true;
            // 对照实验自动选中其他实验组
            const labelSelectedExpItemArray = _.reject(this.labExpItemArray, {
                id: baseExpItemId,
            });
            this.form.lab_experiment_item_id = labelSelectedExpItemArray.map((expItem) => expItem.id);
        },
        initRadarStore(radarId) {
            this.$store.commit('experiment/radar/setStageMap', {
                radarId,
            });
        },
    },
};
</script>
<style lang="scss">
.date-range-expand-checkbox {
    margin-left: 10px;
}
</style>
