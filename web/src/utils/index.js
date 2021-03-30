import DateFormat from 'dateformat';
import _ from 'lodash';

/**
 * 日期格式化方法
 */
export const DateTimeFormat = {
    /**
     * 将时间戳转成对应时间格式
     * @param {Number} timestamp 时间戳
     * @param {String} format 'yyyy-mm-dd HH:MM:ss'
     */
    formatByTimestamp: (timestamp, format) => {
        let result = '';
        const _TAG = 'formatByTimestamp';
        do {
            if (!isNotEmptyString(format)) {
                console.error(`${_TAG} format is invalid string`);
                break;
            }
            if (!_.isNumber(timestamp)) {
                timestamp = parseInt(timestamp, 10);
            }
            if (_.isNaN(timestamp)) {
                console.error(`${_TAG} timestamp is invalid`);
                break;
            }
            try {
                result = DateFormat(new Date(timestamp), format);
            } catch (err) {
                console.error(`${_TAG} ${err.toString()}`);
            }
        } while (false);
        return result;
    },
    formatByDate: (date, format) => {
        let result = '';
        const _TAG = 'formatByDate';
        do {
            if (!isNotEmptyString(format)) {
                console.error(`${_TAG} format is invalid string`);
                break;
            }
            try {
                result = DateFormat(new Date(date), format);
            } catch (err) {
                console.error(`${_TAG} ${err.toString()}`);
            }
        } while (false);
        return result;
    },
};

/**
 * 支持将 id, parent_id 表示的数组转为树形结构
 * [{ id: 1, parent_id: 3 },
 *  { id: 2, parent_id: 0 },
 *  { id: 3, parent_id: 0 },
 *  { id: 4, parent_id: 1 }]
 * ===>
 * [
 *  { id: 2, children: [] },
 *  { id: 3, children: [
 *      { id: 1, children: [
 *          { id: 4, children: [] }
 * ]}]}]
 * @param {Array} data
 * @param {String} options.idKey data中标识id的关键字, 默认 id
 * @param {String} options.parentIdKey data中标识parent_id的关键字, 默认 parent_id
 * @param {String} options.parentIdNullVal data中标识root节点的值, 默认 null
 * @param {String} options.childrenKey 结果数据中标识子节点的关键字, 默认 children
 */
export const arrayToTree = (data, options) => {
    const TAG = 'Utils.arrayToTree';
    const { idKey = 'id', parentIdKey = 'parent_id', parentIdNullVal = null, childrenKey = 'children' } = options || {};
    const idMapping = data.reduce((acc, el, i) => {
        acc[el[idKey]] = i;
        return acc;
    }, {});
    const rootArray = [];
    data.forEach((el) => {
        const parentId = el[parentIdKey];
        if (parentId === undefined || parentId === parentIdNullVal) {
            rootArray.push(el);
            return;
        }
        const parentEl = data[idMapping[parentId]];
        if (parentEl === undefined) {
            console.error(`${TAG}: parent节点为空, 请检查数据或配置`);
            return;
        }
        parentEl[childrenKey] = [...(parentEl[childrenKey] || []), el];
    });
    return rootArray;
};

export const isNotEmptyString = (str) => {
    return isStr(str) && !isEmpty(str.trim());
};

export const isNonEmptyArray = (array) => {
    return isArray(array) && !isEmpty(array);
};

/**
 * 检查对象是否为空, 支持String、Object、Array、Map 和 Set
 * @param {*} x
 */
export const isEmpty = (x) => {
    if (Array.isArray(x) || typeof x === 'string' || x instanceof String) {
        return x.length === 0;
    }

    if (x instanceof Map || x instanceof Set) {
        return x.size === 0;
    }

    if ({}.toString.call(x) === '[object Object]') {
        return Object.keys(x).length === 0;
    }

    return false;
};

const is = function(obj, name) {
    return Object.prototype.toString.call(obj).replace(/(^\[object )|(]$)/g, '') === name;
};

const isArray = (obj) => is(obj, 'Array');
const isObject = (obj) => is(obj, 'Object');
const isNum = (obj) => is(obj, 'Number');
const isStr = (obj) => is(obj, 'String');
const isBool = (obj) => is(obj, 'Boolean');
const isFunction = (obj) => is(obj, 'Function');
const isEvent = (obj) => is(obj, 'Event');
const isNull = (obj) => obj === null;
const isPositiveInteger = (obj) => _.inRange(obj, 0, Number.MAX_SAFE_INTEGER);

const asyncSequentializer = (() => {
    const toPromise = (x) => {
        // if promise just return it
        if (x instanceof Promise) {
            return x;
        }
        // if function is not async this will turn its result into a promise
        // if it is async this will await for the result
        if (typeof x === 'function') {
            return (async () => await x())();
        }

        return Promise.resolve(x);
    };

    return (list) => {
        const results = [];

        return (
            list
                .reduce((lastPromise, currentPromise) => {
                    return lastPromise.then((res) => {
                        results.push(res); // collect the results
                        return toPromise(currentPromise);
                    });
                }, toPromise(list.shift()))
                // collect the final result and return the array of results as resolved promise
                .then((res) => Promise.resolve([...results, res]))
        );
    };
})();

export default {
    DateTimeFormat,
    isNotEmptyString,
    isArray,
    isObject,
    isNum,
    isStr,
    isBool,
    isFunction,
    isEmpty,
    isEvent,
    isNull,
    isPositiveInteger,
    isNonEmptyArray,
    is,
    asyncSequentializer,
    arrayToTree,
};
