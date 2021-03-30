import _ from 'lodash';
import jsonpatch from 'fast-json-patch';
import Utils from '@/utils';

/**
 * Diff 字段配置
 *  * key: 配置唯一key
 *  * path: 用于查找jsonpatch差异的字段路径
 *  * regex: 用于匹配jsonpatch差异字段路径的正则表达式
 *  * values: 用于展示差异结果的function, 默认返回 item.value
 *  * label: 字段名称
 */
export const config = {
    group: [
        {
            key: 'name',
            path: '/data/name',
            label: '实验组名称',
        },
        {
            key: 'description',
            path: '/data/description',
            label: '实验描述',
        },
        {
            key: 'rta_account_id',
            path: '/data/rta_account_id',
            label: 'RTA账号',
        },
        {
            key: 'user',
            regex: /^\/data\/user\/[0-9]+(\/name)?$/,
            values: (items) => {
                const value = items
                    .map((item) => {
                        if (Utils.isNotEmptyString(item.value)) {
                            return item.value;
                        }
                        return item.value.name;
                    })
                    .join(',');
                return value;
            },
            label: '关注人',
        },
    ],
    stage: [
        {
            key: 'end_time',
            path: '/data/current/end_time',
            label: '结束时间',
        },
    ],
    item: [
        {
            key: 'end_time',
            path: 'end_time',
            label: '结束时间',
        },
    ],
};

export const PATCH_TYPE = {
    Modify: {
        key: 'modify',
        label: '修改',
        class: 'warning',
    },
    Create: {
        key: 'create',
        label: '新增',
        class: 'success',
    },
    Delete: {
        key: 'delete',
        label: '删除',
        class: 'danger',
    },
};

const Diff = {
    /**
     * 处理并使用jsonpatch计算record的差异
     * @param {Object} lhsRecord 上一次Modify状态的记录
     * @param {Object} rhsRecord 本次Modify状态的记录
     * @returns {diffResult} 差异数组
     *  [{
     *      "type": PATCH_TYPE.Modify,   // Modify, Create, Delete
     *      "path": "name",
     *      "label": "实验组名称",
     *      "lhs": "ExpGroup - 1",
     *      "rhs": "ExpGroup - 2"
     *  }]
     * @returns {visualizeResult} 展示差异的el-tree节点数据
     *  [{
     *      "label": "基础信息",
     *      "children": [{
     *          "type": PATCH_TYPE.Modify,
     *          "label": "实验名称: lhs → rhs"
     *      }]
     *  }, {
     *      "label": "实验AATest",
     *      "type": "type": PATCH_TYPE.Create
     *  }]
     */
    diffRecords: (lhsRecord, rhsRecord) => {
        const diffResult = {
            base: [],
            item: [],
        };
        // Diff Group层级
        const groupResult = Diff.diffGroup(lhsRecord, rhsRecord);
        if (Utils.isNonEmptyArray(groupResult)) {
            diffResult.base = diffResult.base.concat(groupResult);
        }
        // Diff Stage层级
        const stageResult = Diff.diffStage(lhsRecord, rhsRecord);
        if (Utils.isNonEmptyArray(stageResult)) {
            diffResult.base = diffResult.base.concat(stageResult);
        }
        // Diff Item层级
        const itemResult = Diff.diffItem(lhsRecord, rhsRecord);
        if (Utils.isNonEmptyArray(itemResult)) {
            diffResult.item = diffResult.item.concat(itemResult);
        }
        // 将diff结果转换成el-tree的节点数据
        const visualizeResult = Diff.visualDiffTree(diffResult);
        return {
            diffResult,
            visualizeResult,
        };
    },
    diffGroup: (lhsRecord, rhsRecord) => Diff.patchAndVisualize(lhsRecord, rhsRecord, config.group),
    diffStage: (lhsRecord, rhsRecord) => Diff.patchAndVisualize(lhsRecord, rhsRecord, config.stage),
    diffItem: (lhsRecord, rhsRecord) => {
        const diffResult = [];
        const lhsExpItem = lhsRecord?.data?.current?.experiment_item;
        const rhsExpItem = rhsRecord?.data?.current?.experiment_item;

        // 根据right-hand-side数据筛选出 新增 和 修改
        if (Utils.isNonEmptyArray(rhsExpItem)) {
            rhsExpItem.forEach((rhsItem) => {
                const { outer_id: outerId, name } = rhsItem;
                const lhsItem = _.remove(lhsExpItem, (item) => {
                    return item.outer_id === outerId;
                })[0];
                if (typeof lhsItem === 'undefined') {
                    // 新增
                    diffResult.push({
                        id: outerId,
                        name,
                        type: PATCH_TYPE.Create,
                    });
                } else {
                    // 修改
                    const modifyList = []
                        .concat(Diff.patchAndVisualize(lhsItem, rhsItem, config.item))
                        .concat(Diff.diffRtaExp(lhsItem, rhsItem))
                        .concat(Diff.diffMetadata(lhsItem, rhsItem));
                    if (modifyList.length > 0) {
                        diffResult.push({
                            id: outerId,
                            name,
                            type: PATCH_TYPE.Modify,
                            list: modifyList,
                        });
                    }
                }
            });
        }
        // 根据left-hand-side中剩余的item筛选出 删除
        if (Utils.isNonEmptyArray(lhsExpItem)) {
            // 删除
            lhsExpItem.forEach((lhsItem) => {
                const { outer_id: outerId, name } = lhsItem;
                diffResult.push({
                    id: outerId,
                    name,
                    type: PATCH_TYPE.Delete,
                });
            });
        }
        return diffResult;
    },
    diffRtaExp: (lhsItem, rhsItem) => {
        const lhsRtaExpId = Utils.isNonEmptyArray(lhsItem.rta_exp) ? lhsItem.rta_exp.map((rtaExp) => rtaExp.id) : [];
        const rhsRtaExpId = Utils.isNonEmptyArray(rhsItem.rta_exp) ? rhsItem.rta_exp.map((rtaExp) => rtaExp.id) : [];
        // 使用 _.difference 来兼容数组序号变化带来的diff
        const diff = _.difference(lhsRtaExpId, rhsRtaExpId);
        if (diff.length > 0) {
            return [
                {
                    type: PATCH_TYPE.Modify,
                    key: 'rta_exp',
                    label: 'RTA EXP绑定状态',
                    lhs: lhsRtaExpId.join(','),
                    rhs: rhsRtaExpId.join(','),
                },
            ];
        } else {
            return [];
        }
    },
    diffMetadata: (lhsItem, rhsItem) => {
        const diffResult = [];
        const lhsMetadata = lhsItem.experiment_metadata;
        const rhsMetadata = rhsItem.experiment_metadata;
        // 根据right-hand-side数据筛选出 新增 和 修改
        if (Utils.isNonEmptyArray(rhsMetadata)) {
            rhsMetadata.forEach((rhsMeta) => {
                const {
                    experiment_parameter_id: rhsParamId,
                    experiment_parameter: rhsParam,
                    value: rhsValue,
                } = rhsMeta;
                const lhsMeta = _.remove(lhsMetadata, (meta) => meta.experiment_parameter_id === rhsParamId)[0];
                if (typeof lhsMeta === 'undefined') {
                    // 新增
                    diffResult.push({
                        type: PATCH_TYPE.Create,
                        path: rhsParamId,
                        label: rhsParam.name,
                        rhs: rhsValue,
                    });
                } else {
                    // 修改
                    const { value: lhsValue } = lhsMeta;
                    if (lhsValue !== rhsValue) {
                        diffResult.push({
                            type: PATCH_TYPE.Modify,
                            path: rhsParamId,
                            label: rhsParam.name,
                            lhs: lhsValue,
                            rhs: rhsValue,
                        });
                    }
                }
            });
        }
        // 根据left-hand-side中剩余的item筛选出 删除
        if (Utils.isNonEmptyArray(lhsMetadata)) {
            // 删除
            lhsMetadata.forEach((lhsMeta) => {
                const {
                    experiment_parameter_id: lhsParamId,
                    experiment_parameter: lhsParam,
                    value: lhsValue,
                } = lhsMeta;
                diffResult.push({
                    type: PATCH_TYPE.Delete,
                    path: lhsParamId,
                    label: lhsParam.name,
                    lhs: lhsValue,
                });
            });
        }
        return diffResult;
    },
    /**
     *  {
     *      "type": “modify",   // modify, create, delete
     *      "path": "name",
     *      "label": "实验组名称",
     *      "lhs": "ExpGroup - 1",
     *      "rhs": "ExpGroup - 2"
     *  }
     * op: test, remove, replace, move, add
     */
    patchAndVisualize: (lhs, rhs, config) => {
        const patchResult = [];
        const patch = jsonpatch.compare(lhs, rhs, true);
        if (Utils.isNonEmptyArray(patch)) {
            config.some((conf) => {
                const { path, label } = conf;
                const matchFunc = Diff.getMatchFunc(conf);
                const valueFunc = Diff.getValueFunc(conf);
                if (!Utils.isFunction(matchFunc) || !Utils.isFunction(valueFunc)) {
                    return false;
                }
                const testPatchItem = _.filter(matchFunc(patch), {
                    op: 'test',
                });
                const replacePatchItem = _.filter(matchFunc(patch), {
                    op: 'replace',
                });
                let addPatchItem = _.filter(matchFunc(patch), {
                    op: 'add',
                });
                const removePatchItem = _.filter(matchFunc(patch), {
                    op: 'remove',
                });
                // 新增
                if (addPatchItem.length > 0) {
                    // 新增需要处理因为数组顺序调整而带来的test和replace
                    addPatchItem = _.difference(_.union(addPatchItem, replacePatchItem), testPatchItem);
                    patchResult.push({
                        type: PATCH_TYPE.Create,
                        path,
                        label,
                        rhs: valueFunc(addPatchItem),
                    });
                } else if (testPatchItem.length > 0 && removePatchItem.length > 0) {
                    // 删除
                    patchResult.push({
                        type: PATCH_TYPE.Delete,
                        path,
                        label,
                        lhs: valueFunc(testPatchItem),
                    });
                } else if (testPatchItem.length > 0 && replacePatchItem.length > 0) {
                    // 修改
                    patchResult.push({
                        type: PATCH_TYPE.Modify,
                        path,
                        label,
                        lhs: valueFunc(testPatchItem),
                        rhs: valueFunc(replacePatchItem),
                    });
                } else {
                    return false;
                }
            });
        }
        return patchResult;
    },
    getMatchFunc(conf) {
        const { path, regex } = conf || {};
        if (_.isRegExp(regex)) {
            return (patch) => {
                return _.filter(patch, (o) => {
                    return regex.test(o.path);
                });
            };
        } else if (Utils.isNotEmptyString(path)) {
            return (patch) => {
                return _.filter(patch, {
                    path,
                });
            };
        } else {
            return;
        }
    },
    getValueFunc(conf) {
        const { values } = conf || {};
        if (Utils.isFunction(values)) {
            return (items) => {
                return values(items);
            };
        } else {
            return (items) => {
                return items.map((item) => item.value).join(',');
            };
        }
    },
    visualDiffTree(diffResult) {
        const visualTree = [];
        const { base, item } = diffResult || {};
        if (Utils.isNonEmptyArray(base)) {
            const baseTree = {};
            baseTree.label = '基本信息';
            baseTree.children = [];
            base.forEach((patchItem) => {
                baseTree.children.push(Diff.generatePatchSpan(patchItem));
            });
            visualTree.push(baseTree);
        }
        if (Utils.isNonEmptyArray(item)) {
            item.forEach((itemPatch) => {
                const { name: itemName, type: itemPatchType, list: itemPatchList } = itemPatch || {};
                const { key: patchKey } = itemPatchType || {};
                const itemTree = {};
                switch (patchKey) {
                    case PATCH_TYPE.Modify.key: {
                        itemTree.label = `实验 ${itemName}`;
                        itemTree.patchType = itemPatchType;
                        itemTree.children = [];
                        itemPatchList.forEach((patchItem) => {
                            itemTree.children.push(Diff.generatePatchSpan(patchItem));
                        });
                        break;
                    }
                    case PATCH_TYPE.Create.key:
                    case PATCH_TYPE.Delete.key: {
                        itemTree.label = `实验 ${itemName}`;
                        itemTree.patchType = itemPatchType;
                        break;
                    }
                }
                visualTree.push(itemTree);
            });
        }
        return visualTree;
    },
    generatePatchSpan(patchItem) {
        const { type: patchType, label, lhs, rhs } = patchItem || {};
        const { key: patchKey } = patchType || {};
        let spanLabel = '';
        switch (patchKey) {
            case PATCH_TYPE.Modify.key: {
                spanLabel = `${label}：${lhs || '-'} → ${rhs || '-'}`;
                break;
            }
            case PATCH_TYPE.Create.key: {
                spanLabel = `${label}：${rhs}`;
                break;
            }
            case PATCH_TYPE.Delete.key: {
                spanLabel = `${label}：${lhs}`;
                break;
            }
        }
        return {
            patchType,
            label: spanLabel,
        };
    },
};
export default Diff;
