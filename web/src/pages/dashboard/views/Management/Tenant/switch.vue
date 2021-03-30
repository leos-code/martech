<template>
    <el-dialog title="选择组织" :visible.sync="dialogVisible" @close="close">
        <el-table
            v-if="tenant.length > 0"
            :ref="tableRef"
            :data="tenant"
            highlight-current-row
            @current-change="handleTenantSelected"
        >
            <el-table-column label="选择">
                <template slot-scope="scope">
                    <el-radio v-model="selectedTenant" :label="scope.row">
                        <i />
                    </el-radio>
                </template>
            </el-table-column>
            <el-table-column label="组织名称" property="name" />
        </el-table>
        <span v-if="tenant.length > 0" slot="footer">
            <el-button @click="close">
                取 消
            </el-button>
            <el-button type="primary" @click="commit">
                确 认
            </el-button>
        </span>
        <div v-else>
            您当前未加入任何组织, 请联系组织管理员邀请加入
        </div>
    </el-dialog>
</template>
<script>
import _ from 'lodash';
import { mapGetters } from 'vuex';
import { User } from '@dashboard/api';
import Utils from '@/utils';
import Casbin from '@dashboard/Casbin';

const TAG = 'TenantSwitch';

export default {
    name: 'TenantSwitch',
    data() {
        return {
            tableRef: 'tenantSwitchTableRef',
            dialogVisible: false,
            selectedTenant: undefined,
        };
    },
    computed: {
        ...mapGetters('user', ['currentTenant', 'tenant', 'tenantSwitchActive']),
    },
    watch: {
        currentTenant: {
            handler() {
                this.init();
            },
            deep: true,
        },
        tenantSwitchActive(active) {
            this.dialogVisible = active;
        },
    },
    mounted() {
        this.init();
    },
    methods: {
        show() {
            this.dialogVisible = true;
            this.$store.dispatch('user/activeTenantSwitch');
        },
        close() {
            this.dialogVisible = false;
            this.$store.dispatch('user/closeTenantSwitch');
        },
        async init() {
            await this.$store.dispatch('user/initUserInfo');
            const { id: currentTenantId } = this.currentTenant;
            if (!Utils.isPositiveInteger(currentTenantId)) {
                this.show();
            } else {
                this.selectedTenant = _.find(this.tenant, { id: currentTenantId });
            }
        },
        handleTenantSelected(row) {
            this.selectedTenant = row;
        },
        commit() {
            const msgPrefix = '切换组织';
            if (typeof this.selectedTenant === undefined) {
                this.$message.warning(`${msgPrefix}失败 - 请先选择一个组织`);
                return false;
            }
            if (this.selectedTenant?.id === this.currentTenant?.id) {
                this.$message.warning(`${msgPrefix}失败 - 您当前已处于该组织，无需切换`);
                return false;
            }
            User.tenant(this.selectedTenant)
                .then(async () => {
                    this.$message.success(`${msgPrefix}成功 - 马上为您跳转`);
                    await Casbin.clear();
                    this.$router.go(0);
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                });
        },
    },
};
</script>
