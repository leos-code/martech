import _ from 'lodash';
import Utils from '@/utils';

export const getSpecStatusStage = (group, status) => {
    if (typeof group === 'undefined') {
        return [];
    }
    const stage = _.get(group, 'experiment_stage');
    const specStatusStage = _.filter(stage, {
        status,
    });
    return specStatusStage;
};

export const DataFormat = {
    format: (value, formatter) => {
        if (typeof formatter === 'undefined') {
            return value;
        } else {
            let formatArr = [];
            if (Utils.isNotEmptyString(formatter)) {
                formatArr.push({
                    type: formatter,
                });
            } else if (Utils.isObject(formatter)) {
                formatArr.push(formatter);
            } else if (Utils.isArray(formatter)) {
                formatArr = formatter;
            } else if (Utils.isFunction(formatter)) {
                formatArr.push({
                    func: formatter,
                });
            }
            formatArr.map((item) => {
                const { type, config, func } = item;
                if (Utils.isFunction(func)) {
                    value = func(value, config);
                } else if (type === 'percentage') {
                    const { precision = 2, defaultValue } = config || {};
                    value = DataFormat.formatPercentage(value, precision, defaultValue);
                } else if (type === 'thousandSeparator') {
                    const { precision } = config || {};
                    value = DataFormat.formatThousandsSpace(value, precision);
                } else if (type === 'round') {
                    const { precision = 3 } = config || {};
                    value = _.round(value, precision);
                } else if (type === 'date') {
                    console.log('date');
                }
            });
            return value;
        }
    },
    // 转换为百分比
    formatPercentage: (value, precision, defaultValue) => {
        // 小数点后 precision 位，四舍五入
        const v = _.round((parseFloat(value) || 0) * 100, precision);
        if (v == 0 && defaultValue) {
            value = String(defaultValue);
        } else {
            value = `${v}%`;
        }
        return value;
    },
    /**
     * 数字千分位转换显示
     */
    formatThousandsSpace: (num, precision) => {
        if (_.isNumber(precision) && _.isNumber(num)) {
            num = num.toFixed(precision);
        }
        const str = String(num);
        const idx = str.indexOf('.');
        const intPart = idx > -1 ? str.substring(0, idx) : str;
        const floatPart = idx > -1 ? str.substring(idx) : '';
        const reg = /\d{1,3}(?=(\d{3})+$)/g;
        return intPart.replace(reg, '$&,') + floatPart;
    },
};
