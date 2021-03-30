import Casbin from '@dashboard/Casbin';
import Utils from '@/utils';

/**
 * 页面内功能按钮的权限控制，组件上通过 v-permission="'subFunctionKey'" 配置
 */
async function checkPermission(el, binding) {
    const { value } = binding || {};

    if (Utils.isNotEmptyString(value)) {
        const authResult = await Casbin.can({
            object: `${value}#sub_function`,
            action: 'read',
        });
        if (!authResult) {
            el.parentNode && el.parentNode.removeChild(el);
        }
    } else {
        throw new Error(`Need permissions! Like v-permission="materialAudit"`);
    }
}

export default {
    async inserted(el, binding) {
        await checkPermission(el, binding);
    },
    async update(el, binding) {
        await checkPermission(el, binding);
    },
};
