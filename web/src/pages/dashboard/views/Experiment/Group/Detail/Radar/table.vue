<template>
    <el-table
        border
        class="experiment-radar-table"
        :data="tableData"
        :cell-class-name="setCellClassName"
        :span-method="spanMethod"
    >
        <!-- 维度 -->
        <el-table-column v-for="dimention in Dimentions" :key="dimention.key" :label="dimention.name">
            <template slot-scope="scope">
                {{ getColValue(scope.row, dimention.key) }}
            </template>
        </el-table-column>
        <!-- 指标 -->
        <el-table-column v-for="field in Fields" :key="field.key">
            <template slot="header">
                <span>{{ field.name }}</span>
                <el-tooltip v-if="isNotEmptyString(field.tips)" effect="dark" :content="field.tips" placement="top">
                    <i class="el-icon-question" />
                </el-tooltip>
            </template>
            <template slot-scope="scope">
                <div>{{ getColValue(scope.row, field.key).value }}</div>
                <div v-if="!scope.row.is_base" :class="getColValue(scope.row, field.key).color">
                    {{ getColValue(scope.row, field.key).delta }}
                </div>
            </template>
        </el-table-column>
    </el-table>
</template>
<script>
import _ from 'lodash';
import Utils from '@/utils';
import { DataFormat } from '../utils';
import { Dimentions, Fields, Granularity } from './meta';
import { ExperimentRadar } from '@dashboard/api';

const baseFieldConfig = {
    key: 'key',
    name: '指标名称',
};

export default {
    name: 'ExperimentGroupDetailRadarTable',
    props: {
        group: Object,
        stage: Object,
    },
    data() {
        return {
            Dimentions,
            Fields,
            tableData: [],
        };
    },
    computed: {
        radarId() {
            return `radar-group${this.group.id}-stage${this.stage.id}`;
        },
        getQuerying() {
            return this.$store.getters['experiment/radar/getQuerying'](this.radarId);
        },
        queryOptions() {
            return this.$store.getters['experiment/radar/queryOptions'](this.radarId);
        },
    },
    watch: {
        getQuerying(querying) {
            if (querying) {
                this.getRadar(this.queryOptions);
                this.$store.commit('experiment/radar/closeQuerying', {
                    radarId: this.radarId,
                });
            }
        },
        queryOptions() {},
    },
    methods: {
        isNotEmptyString: Utils.isNotEmptyString,
        getRadar(options) {
            const msgPrefix = '查询效果数据';
            ExperimentRadar.get(options)
                .then(({ data }) => {
                    const { records } = data || {};
                    this.tableData = this.compileTableData(records);
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix} - ${msg}`);
                });
        },
        compileTableData(data) {
            if (!Utils.isNonEmptyArray(data)) {
                return [];
            }
            const resultDataArray = _.cloneDeep(data);
            // 维度&指标 格式化
            const columnConfigs = [].concat(this.Dimentions).concat(this.Fields);
            const { base_experiment_item_id: baseExpItemId, granularity = 'day' } = this.queryOptions;
            resultDataArray.map((row) => {
                // 判断是否为基准实验组
                row.is_base = row.experiment_item.id === baseExpItemId;
                // 根据粒度和时间格式，格式化展示的时间
                row.time = this.timestampFormat({
                    timestampArr: _.get(row, 'time'),
                    granularityKey: granularity,
                });
                columnConfigs.map((columnConfig) => {
                    if (typeof columnConfig.values === 'undefined') {
                        this.compileCellData(row, columnConfig);
                    } else {
                        const { key: prefixKey, values } = columnConfig;
                        values.map((valueConfig) => {
                            this.compileCellData(row, valueConfig, prefixKey);
                        });
                    }
                });
            });
            return resultDataArray;
        },
        compileCellData(row, columnConfig = {}, prefixKey = '') {
            if (!_.isObject(row) || !_.isObject(columnConfig)) {
                return;
            }
            const { key: childKey, formatter, colorful } = columnConfig || {};
            const fullPath = Utils.isNotEmptyString(prefixKey) ? `${prefixKey}.${childKey}` : childKey;
            const originValue = _.get(row, fullPath);
            // 格式化展示字段
            const showValue = DataFormat.format(originValue, formatter);
            _.set(row, fullPath, showValue);
            // 格式化正负值
            let color = '';
            if (colorful && _.isNumber(originValue) && Utils.isNotEmptyString(prefixKey)) {
                if (originValue > 0) {
                    color = 'positive';
                } else if (originValue < 0) {
                    color = 'negative';
                }
                _.set(row, `${prefixKey}.color`, color);
            }
        },
        getColValue(row, key) {
            return _.get(row, key, {});
        },
        mergeFieldConfig(config) {
            return _.assign({}, baseFieldConfig, config);
        },
        getSpanArray(data) {
            const spanArray = [];
            const dimentionLen = this.Dimentions.length;
            if (dimentionLen < 2) {
                return spanArray;
            }
            _.reverse(this.Dimentions.slice(1, dimentionLen)).map((dimention, index) => {
                const rowSpan = _.uniq(data.map((item) => _.get(item, dimention.key))).length;
                spanArray.push({
                    rowSpan,
                    columnIndex: index,
                });
            });
            return spanArray;
        },
        spanMethod({ rowIndex, columnIndex }) {
            let rowspan = 1;
            let colspan = 1;
            const spanArray = this.getSpanArray(this.tableData);
            spanArray.map((spanItem) => {
                const { rowSpan: spanedRowSpan, columnIndex: spanedColIndex } = spanItem || {};
                if (columnIndex === spanedColIndex) {
                    if (rowIndex % spanedRowSpan === 0) {
                        rowspan = spanedRowSpan;
                        colspan = 1;
                    } else {
                        rowspan = 0;
                        colspan = 0;
                    }
                }
            });
            return {
                rowspan,
                colspan,
            };
        },
        setCellClassName({ row }) {
            const { is_base } = row || {};
            return is_base ? 'base-row' : '';
        },
        /**
         * 时间戳存在两种格式:
         *  - 时间点, 如 [1607731200000, 0] => 2020-12-12
         *  - 时间段, 如 [1607731200000, 1610582400000] => [2020-12-12, 2021-01-14]
         */
        timestampFormat({ timestampArr, granularityKey }) {
            if (!Utils.isNonEmptyArray(timestampArr)) {
                return [];
            }
            const granularity = _.find(Granularity, {
                key: granularityKey,
            });
            const resultTimeArr = timestampArr
                .filter((timestamp) => timestamp !== 0)
                .map((timestamp) => {
                    return Utils.DateTimeFormat.formatByTimestamp(timestamp * 1000, granularity.dateFmt);
                });
            return resultTimeArr.join('-');
        },
    },
};
</script>
<style lang="scss">
.experiment-radar-table {
    margin-top: 40px;
    .base-row {
        background-color: #fafafa;
    }
    td {
        padding: 4px 0;
        .cell {
            text-align: center;
            font-size: 13px;
        }
    }
    th > .cell {
        text-align: center;
    }
    .positive {
        color: #5cbb7a;
    }
    .negative {
        color: #f56c6c;
    }
}
</style>
