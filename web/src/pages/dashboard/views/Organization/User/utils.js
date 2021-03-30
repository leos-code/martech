import _ from 'lodash';
import Utils from '@/utils';
import { LOGIN_TYPE } from './meta';

/**
 * 按照priority筛选loginUser
 * @param {Array} data login user 列表
 * @returns {Object} 最高优先级的login user
 */
export const getVisibleLoginUser = (data) => {
    if (!Utils.isNonEmptyArray(data)) {
        return {};
    }
    let loginUserList = _.cloneDeep(data);
    loginUserList.forEach((loginUser) => {
        const { login_type } = loginUser || {};
        loginUser.priority = LOGIN_TYPE[login_type].priority || 999;
    });
    loginUserList = _.sortBy(loginUserList, ['priority']);
    return loginUserList[0];
};
