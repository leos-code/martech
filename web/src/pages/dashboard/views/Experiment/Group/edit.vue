<template>
    <el-row :gutter="24">
        <el-col :span="18" :offset="2">
            <el-form :ref="formRef" :model="form" label-width="80px" :rules="rules.group">
                <el-divider content-position="left">
                    基本信息
                </el-divider>
                <el-form-item label="实验组名称" label-width="160px" prop="name">
                    <el-input v-model="form.name" autocomplete="off" />
                </el-form-item>
                <el-form-item label="实验目标" label-width="160px">
                    <el-input
                        v-model="form.description"
                        type="textarea"
                        :rows="3"
                        placeholder="请输入实验目标或描述信息"
                    />
                </el-form-item>
                <el-form-item label="RTA账户" label-width="160px" prop="rta_account_id">
                    <el-select
                        v-model="form.rta_account_id"
                        placeholder="请选择RTA账户"
                        disabled
                        @change="changeRtaAccount"
                    >
                        <el-option
                            v-for="account in rtaAccounts"
                            :key="account.id"
                            :label="account.rta_id"
                            :value="account.id"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item label="关注人" label-width="160px" prop="user_id">
                    <el-select
                        v-model="form.user_id"
                        placeholder="请选择关注人"
                        multiple
                        filterable
                        :default-first-option="true"
                        class="select-user-id"
                        @change="changeUserId"
                    >
                        <el-option v-for="user in Users" :key="user.id" :label="user.name" :value="user.id" />
                    </el-select>
                </el-form-item>
                <el-form-item label="结束时间" label-width="160px" prop="draft.end_time">
                    <el-date-picker v-model="form.draft.end_time" type="datetime" placeholder="选择结束时间" />
                </el-form-item>
                <el-divider content-position="left">
                    实验管理
                </el-divider>
                <el-button type="primary" size="small" @click="addItem">
                    新增实验
                </el-button>
                <el-tabs
                    v-if="form.draft && form.draft.experiment_item && form.draft.experiment_item.length > 0"
                    v-model="activeItem"
                    type="card"
                    closable
                    class="experiment-item-tabs"
                    @tab-remove="removeItem"
                >
                    <el-tab-pane
                        v-for="(exp_item, exp_item_index) in form.draft.experiment_item"
                        :key="exp_item_index"
                        :label="exp_item.name || `实验-${exp_item_index}`"
                        :name="exp_item_index + ''"
                    >
                        <el-form-item label="实验名称" label-width="160px">
                            <el-input v-model="exp_item.name" autocomplete="off" />
                        </el-form-item>
                        <el-form-item label="RtaExpID" label-width="160px">
                            <el-select
                                v-model="exp_item.selected_rta_exp_ids"
                                multiple
                                placeholder="请选择"
                                @change="changeRtaExp(exp_item.selected_rta_exp_ids, exp_item_index)"
                            >
                                <el-option
                                    v-for="rtaExpId in rtaExpIds"
                                    :key="rtaExpId.id"
                                    :label="`${rtaExpId.id} - ${rtaExpId.flow_rate}%`"
                                    :value="rtaExpId.id"
                                    :disabled="
                                        exp_item.disable_rta_exp_ids &&
                                            exp_item.disable_rta_exp_ids.includes(rtaExpId.id)
                                    "
                                />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="采样比例" label-width="160px" class="flow-rate-item">
                            <span>{{ exp_item.total_flow_rate }}%</span>
                        </el-form-item>
                        <el-form-item label="实验参数" label-width="160px">
                            <el-input
                                v-for="(exp_metadata, exp_metadata_index) in exp_item.experiment_metadata"
                                :key="exp_metadata_index"
                                v-model="exp_metadata.value"
                                placeholder="请输入参数值"
                                class="input-with-select"
                            >
                                <el-select
                                    slot="prepend"
                                    v-model="exp_metadata.experiment_parameter_id"
                                    placeholder="请选择参数"
                                >
                                    <el-option
                                        v-for="parameter in experimentParameters"
                                        :key="parameter.id"
                                        :label="parameter.name"
                                        :value="parameter.id"
                                    />
                                </el-select>
                                <el-button
                                    slot="append"
                                    icon="el-icon-delete"
                                    @click="
                                        removeMetadata(exp_item.experiment_metadata, exp_metadata_index, exp_item_index)
                                    "
                                />
                            </el-input>
                            <el-button size="small" type="primary" @click="addMetadata(exp_item, exp_item_index)">
                                添加参数
                            </el-button>
                        </el-form-item>
                    </el-tab-pane>
                </el-tabs>
                <el-divider />
                <el-form-item>
                    <el-button type="primary" @click="onSubmit">
                        提交
                    </el-button>
                    <el-button @click="onCancel">
                        取消
                    </el-button>
                </el-form-item>
            </el-form>
        </el-col>
    </el-row>
</template>
<script>
import _ from 'lodash';
import { ExperimentAccount, ExperimentGroup, ExperimentParameter, OrganizationUser } from '@dashboard/api';
import { RtaExpBindStatus, RtaExpStatus } from '@dashboard/views/Experiment/meta';
import { RULES } from './meta';
import Utils from '@/utils';

const DEFAULT_FORM = {
    name: '',
    description: '',
    rta_account_id: undefined,
    user_id: [],
    draft: {
        end_time: undefined,
        experiment_item: [],
    },
};

const DEFAULT_ITEM = {
    name: '',
    rta_exp: [],
    selected_rta_exp_ids: [],
    disable_rta_exp_ids: [],
    total_flow_rate: 0,
    experiment_metadata: [],
};

const DEFAULT_METADATA = {
    experiment_parameter_id: null,
    value: null,
};

export default {
    name: 'ExperimentGroupEdit',
    data() {
        return {
            formRef: 'experimentGroupEditForm',
            form: _.cloneDeep(DEFAULT_FORM),
            rules: RULES,
            activeItem: '0',
            rtaAccounts: [],
            rtaExpIds: [],
            experimentParameters: [],
            Users: [],
        };
    },
    async mounted() {
        Utils.asyncSequentializer([
            this.getUser,
            this.getParameterList,
            this.initForm,
            this.getRtaAccount,
            this.initUserList,
        ]);
    },
    methods: {
        getRtaAccount() {
            const msgPrefix = '拉取RTA账户信息';
            ExperimentAccount.getList()
                .then(({ data }) => {
                    this.rtaAccounts = data;
                    // 写死逻辑，默认使用第一个RTA账户，后期再支持可编辑
                    this.form.rta_account_id = data[0]['id'];
                    this.getRtaExpIds(this.form.rta_account_id);
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                    this.rtaAccounts = [];
                });
        },
        getRtaExpIds(rtaAccountId) {
            const msgPrefix = '拉取RTA EXP信息';
            const account = _.find(this.rtaAccounts, {
                id: rtaAccountId,
            });
            ExperimentAccount.getRtaExpList(account)
                .then(({ data }) => {
                    let { rta_exp } = data || {};
                    this.getEditableRtaExp(rta_exp);
                })
                .catch(({ msg }) => {
                    this.$message.error(`${msgPrefix} - ${msg}`);
                    this.rtaExpIds = [];
                });
        },
        getParameterList() {
            return new Promise((resolve, reject) => {
                const msgPrefix = '拉取实验参数信息';
                ExperimentParameter.getList()
                    .then(({ data }) => {
                        this.experimentParameters = data;
                        resolve();
                    })
                    .catch(({ msg }) => {
                        this.$message.error(`${msgPrefix} - ${msg}`);
                        this.experimentParameters = [];
                        reject();
                    });
            });
        },
        getUser() {
            return new Promise((resolve, reject) => {
                const msgPrefix = '拉取用户信息';
                OrganizationUser.get()
                    .then(({ data }) => {
                        this.Users = data;
                        resolve();
                    })
                    .catch(({ msg }) => {
                        this.$message.error(`${msgPrefix}失败 - ${msg}`);
                        reject();
                    });
            });
        },
        getAvaiableRtaExpList(rtaExpList) {
            return _.filter(rtaExpList, {
                status: RtaExpStatus.Valid.key,
                bind_status: RtaExpBindStatus.Idle.key,
            });
        },
        getEditableRtaExp(globalRtaExp) {
            // 获取当前实验Group可操作的rta_exp列表
            let editableRtaExpList = [];
            // 读取空闲状态的rta_exp
            if (Utils.isNonEmptyArray(globalRtaExp)) {
                editableRtaExpList = editableRtaExpList.concat(this.getAvaiableRtaExpList(globalRtaExp));
            }
            // 读取当前实验Group占用的rta_exp
            const { current } = this.form || {};
            const { experiment_item } = current || {};
            if (Utils.isNonEmptyArray(experiment_item)) {
                experiment_item.forEach((exp_item) => {
                    const { rta_exp: bindRtaExp } = exp_item || {};
                    editableRtaExpList = Utils.isNonEmptyArray(bindRtaExp)
                        ? editableRtaExpList.concat(bindRtaExp)
                        : editableRtaExpList;
                });
            }
            // 空闲状态 + 当前实验占用状态 的rta_exp为可操作列表
            this.rtaExpIds = editableRtaExpList;
            this.updateDisableRtaExp();
        },
        initForm() {
            return new Promise((resolve, reject) => {
                let groupId = this.$route.params.id;
                groupId = parseInt(groupId, 10);
                // ID为0或空时，认为是新建态
                if (groupId) {
                    this.initEditForm(groupId)
                        .then(() => {
                            resolve();
                        })
                        .catch(() => {
                            reject();
                        });
                } else {
                    resolve();
                }
            });
        },
        initEditForm(groupId) {
            return new Promise((resolve, reject) => {
                // 编辑态, 先拉取已保存的实验信息
                const msgPrefix = '拉取实验信息';
                ExperimentGroup.get(groupId)
                    .then(({ data: groupData }) => {
                        this.form = _.assign({}, DEFAULT_FORM, groupData);
                        if (this.form && this.form.draft && Utils.isNonEmptyArray(this.form.draft.experiment_item)) {
                            this.form.draft.experiment_item.map((exp_item, exp_item_index) => {
                                exp_item.total_flow_rate = 0;
                                const selected_rta_exp_ids = [];
                                const { rta_exp, experiment_metadata } = exp_item || {};
                                if (Utils.isNonEmptyArray(rta_exp)) {
                                    rta_exp.map((rta_exp_item) => {
                                        exp_item.total_flow_rate += rta_exp_item.flow_rate;
                                        selected_rta_exp_ids.push(rta_exp_item.id);
                                    });
                                }
                                // 将selected_rta_exp_ids设置为响应式，规避UI视图不更新问题
                                this.$set(
                                    this.form.draft.experiment_item[exp_item_index],
                                    'selected_rta_exp_ids',
                                    selected_rta_exp_ids,
                                );
                                // 删除metadata中experiment_parameter结构体，避免修改数据后，后台读id和结构体发生冲突
                                if (Utils.isNonEmptyArray(experiment_metadata)) {
                                    experiment_metadata.forEach((metadata) => {
                                        delete metadata.experiment_parameter;
                                    });
                                }
                            });
                        }
                        resolve();
                    })
                    .catch(({ err, msg }) => {
                        console.error(err);
                        this.$message.error(`${msgPrefix} - ${msg}`);
                        reject();
                    });
            });
        },
        initUserList() {
            const curUserName = this.$cookies.get('user');
            let curUser;
            if (Utils.isNonEmptyArray(this.Users) && Utils.isNotEmptyString(curUserName)) {
                curUser = _.find(this.Users, {
                    name: curUserName,
                });
            }
            if (typeof curUser === 'undefined') {
                this.$message.error('读取当前用户信息失败, 请刷新重试');
                return;
            }
            const { user } = this.form || {};
            if (Utils.isNonEmptyArray(user)) {
                this.form.user_id = user.map((userItem) => userItem.id);
            } else {
                this.form.user_id = [curUser.id];
                this.form.user = [curUser];
            }
        },
        addItem() {
            if (!Utils.isNonEmptyArray(this.form.draft.experiment_item)) {
                this.$set(this.form.draft, 'experiment_item', []);
            }
            this.form.draft.experiment_item.push(_.cloneDeep(DEFAULT_ITEM));
            this.activeItem = this.getActiveItem();
            this.updateDisableRtaExp();
        },
        removeItem(itemIndex) {
            itemIndex = parseInt(itemIndex, 10);
            this.form.draft.experiment_item = _.filter(this.form.draft.experiment_item, (item, index) => {
                return index != itemIndex;
            });
            this.activeItem = this.getActiveItem();
        },
        getActiveItem() {
            if (this.form && this.form.draft && Utils.isNonEmptyArray(this.form.draft.experiment_item)) {
                return this.form.draft.experiment_item.length - 1 + '';
            } else {
                return '0';
            }
        },
        changeRtaAccount(target) {
            this.getRtaExpIds(target);
        },
        changeRtaExp(selectedRtaExpIds, expItemIndex) {
            const editExperimentItem = this.form.draft.experiment_item[expItemIndex];
            editExperimentItem.rta_exp = [];
            editExperimentItem.total_flow_rate = 0;
            selectedRtaExpIds.map((selectedRtaExpId) => {
                const matchedRtaExp = _.find(this.rtaExpIds, (o) => {
                    return o.id === selectedRtaExpId;
                });
                if (matchedRtaExp) {
                    editExperimentItem.total_flow_rate += matchedRtaExp.flow_rate;
                    editExperimentItem.rta_exp.push(matchedRtaExp);
                }
            });
            this.updateDisableRtaExp();
        },
        updateDisableRtaExp() {
            const totalExpItemList = this.form.draft.experiment_item;
            if (!Utils.isNonEmptyArray(totalExpItemList)) {
                return;
            }
            try {
                totalExpItemList.forEach((curExpItem, curExpIndex) => {
                    curExpItem.disable_rta_exp_ids = totalExpItemList
                        .filter((item, index) => curExpIndex !== index)
                        .map((item) => item.rta_exp)
                        .flat()
                        .filter((item) => typeof item !== 'undefined')
                        .map((item) => item.id);
                });
            } catch (err) {
                console.error(err);
            }
        },
        addMetadata(exp_item) {
            if (!Utils.isNonEmptyArray(exp_item.experiment_metadata)) {
                this.$set(exp_item, 'experiment_metadata', []);
            }
            exp_item.experiment_metadata.push(_.cloneDeep(DEFAULT_METADATA));
        },
        removeMetadata(metadata, remove_metadata_index, item_index) {
            metadata = _.filter(metadata, (item, i) => {
                return i !== remove_metadata_index;
            });
            this.form.draft.experiment_item[item_index].experiment_metadata = metadata;
        },
        changeUserId(userIdList) {
            this.form.user = [];
            if (Utils.isNonEmptyArray(userIdList)) {
                userIdList.forEach((userId) => {
                    const user = _.find(this.Users, {
                        id: userId,
                    });
                    if (user) {
                        this.form.user.push(user);
                    }
                });
            }
        },
        onSubmit() {
            this.$refs[this.formRef].validate(async (valid) => {
                if (valid) {
                    const msgPrefix = '实验创建';
                    ExperimentGroup.edit(this.form)
                        .then(({ data }) => {
                            this.$router.push({
                                name: 'ExperimentGroupDetail',
                                params: {
                                    id: data.id,
                                },
                                query: {
                                    tab: 'draft',
                                },
                            });
                            this.$message.success(`${msgPrefix}成功`);
                        })
                        .catch(({ err, msg }) => {
                            console.error(err);
                            this.$message.error(`${msgPrefix} - ${msg}`);
                        });
                } else {
                    return false;
                }
            });
        },
        onCancel() {
            this.$router.go(-1);
        },
    },
};
</script>
<style lang="scss">
.input-with-select {
    padding-bottom: 10px;
    .el-input-group__prepend {
        background-color: #ffffff;
    }
    .el-select {
        min-width: 200px;
    }
}
.select-user-id {
    min-width: 400px;
}
.experiment-item-tabs {
    margin-top: 10px;
}
.flow-rate-item {
    color: #606266;
}
</style>
