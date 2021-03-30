<template>
    <div>
        <el-form-item label="手机号" :label-width="labelWidth">
            <el-input placeholder="请输入手机号" v-model="phoneNumber">
                <template slot="append">
                    <el-button type="primary" :loading="loading" @click="search">
                        验证
                    </el-button>
                </template>
            </el-input>
        </el-form-item>
        <el-form-item v-if="userList.length > 0">
            <el-table
                :ref="tableRef"
                :data="userList"
                highlight-current-row
                @current-change="handleUserSelected"
                class="search-user-table"
            >
                <el-table-column label="选择">
                    <template slot-scope="scope">
                        <el-radio v-model="selectedRadio" :label="scope.row.id">
                            <i />
                        </el-radio>
                    </template>
                </el-table-column>
                <el-table-column label="手机号" property="phone_number" />
                <el-table-column label="邮箱" property="email" />
                <el-table-column label="绑定账号">
                    <template slot-scope="scope">
                        <el-popover
                            v-for="login_user in scope.row.login_user"
                            :key="login_user.open_id"
                            trigger="hover"
                            placement="top"
                        >
                            <p>账号类型: {{ login_user.login_type_format }}</p>
                            <p>昵称: {{ login_user.nick_name }}</p>
                            <img slot="reference" :src="login_user.avatar" class="user-search-avatar" />
                        </el-popover>
                    </template>
                </el-table-column>
            </el-table>
        </el-form-item>
    </div>
</template>
<script>
import _ from 'lodash';
import { User } from '@dashboard/api';
import Utils from '@/utils';
import { LOGIN_TYPE } from './meta';

const getDefaultData = () => {
    return {
        phoneNumber: '',
        userList: [],
        loading: false,
        tableRef: 'searchUserTable',
        selectedRadio: null,
    };
};

export default {
    name: 'UserSearch',
    props: {
        labelWidth: {
            type: String,
            default: '80px',
        },
    },
    data: getDefaultData,
    methods: {
        resetFields() {
            Object.assign(this.$data, getDefaultData());
        },
        search() {
            const msgPrefix = '验证用户信息';
            if (Utils.isNotEmptyString(this.phoneNumber)) {
                this.loading = true;
                User.search({
                    phone_number: this.phoneNumber,
                })
                    .then(({ data }) => {
                        this.userList = this.format(_.cloneDeep(data));
                        this.loading = false;
                    })
                    .catch(() => {
                        this.$message.error(`${msgPrefix}失败 - 请输入有效用户绑定的手机号`);
                        this.loading = false;
                    });
            } else {
                this.$message.warning(`${msgPrefix}失败 - 请输入手机号`);
            }
        },
        format(data) {
            data.forEach((item) => {
                const { login_user } = item || {};
                login_user.forEach((user) => {
                    user.login_type_format = LOGIN_TYPE[user.login_type]?.label || '未知';
                });
            });
            return data;
        },
        handleUserSelected(val) {
            this.selectedRadio = val?.id;
            this.$emit('select', val);
        },
    },
};
</script>
<style lang="scss">
.user-search-avatar {
    width: 38px;
    height: 38px;
    border-radius: 50%;
}
</style>
