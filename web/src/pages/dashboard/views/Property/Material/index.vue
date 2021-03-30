<template>
    <div class="material_page">
        <div class="ml10 mr10">
            <el-input
                v-model="query.name"
                placeholder="输入素材名称查询"
                size="small"
                class="w200 mr15"
                prefix-icon="el-icon-search"
                @change="getMaterialList(1)"
            />
            <!--            <el-select v-model="query.type" placeholder="全部类型" size="small" class="w200 mr15" @change="getMaterialList(1)">-->
            <!--                <el-option v-for="item in constData.type" :value="item.value" :key="item.value" :label="item.label"></el-option>-->
            <!--            </el-select>-->
            <el-select
                v-model="query.audit_status"
                placeholder="全部审核状态"
                size="small"
                class="w200"
                @change="getMaterialList(1)"
            >
                <el-option
                    v-for="item in constData.auditStatus"
                    :value="item.value"
                    :key="item.value"
                    :label="item.label"
                />
            </el-select>
            <el-button size="small ml15" @click="batchEdit">
                <i class="el-icon-edit-outline mr10" />批量管理
            </el-button>

            <el-button type="primary" class="material_upload_btn" @click="openUploadDetail" size="small" v-permission="'MaterialEdit'">
                <i class="el-icon-menu mr10" />上传详情列表
            </el-button>
            <el-upload
                class="material_upload_btn"
                action="/dashboard/material/file/upload"
                multiple
                :on-success="uploadSuccess"
                :on-error="uploadError"
                :before-upload="beforeUpload"
                :auto-upload="true"
                :show-file-list="false"
                v-permission="'MaterialEdit'"
            >
                <el-button class="upload_btn" type="success" size="small">
                    <i class="el-icon-upload mr10" />批量上传素材
                </el-button>
            </el-upload>
        </div>

        <el-divider />

        <div class="material_handle_bar" v-show="batchEditing">
            已选择
            <span style="color:#409EFF">{{ checkedCount }}</span> 项（<span class="ml5 mr5">
                <el-checkbox :indeterminate="indeterminate" v-model="checkAll" @change="handleCheckAllChange">
                    全选
                </el-checkbox> </span
            >）
            <span class="ml15 mr20">|</span>
            <el-button type="text" @click="batchDelete" v-permission="'MaterialEdit'">批量删除</el-button>
            <el-button type="text" @click="batchReject" v-permission="'MaterialAudit'">批量驳回</el-button>
            <el-button type="text" @click="batchPass" v-permission="'MaterialAudit'">批量通过</el-button>
            <el-button type="text" class="material_handle_close" @click="cancelEdit"><i class="el-icon-close"></i></el-button>
        </div>

        <div v-loading="loading">
            <div v-show="!list.length && !loading">
                <div class="gray material_notice">未上传任何素材。</div>
            </div>
            <div v-show="!list.length && loading">
                <div class="gray material_notice">正在拉取...</div>
            </div>
            <div v-show="!!list.length">
                <template v-for="(item, idx) in listGroup">
                    <el-row :key="idx">
                        <el-col :span="24 / eachLineCount" v-for="(subItem, subIdx) in item" :key="subIdx">
                            <div
                                :class="[
                                    'material_wrap',
                                    !!checkBoxItem[subIdx + eachLineCount * idx] ? 'material_wrap_checked' : '',
                                ]"
                            >
                                <div class="material_inner_wrap">
                                    <div @click="openPreview(subItem)">
                                        <div class="material_detail_wrap" v-if="subItem.data.type == 'image'">
                                            <el-image
                                                style="width: 100%;height: 100%;"
                                                :src="subItem.data.image.url"
                                                fit="scale-down"
                                            >
                                                <div slot="error" class="image-slot">
                                                    <i class="el-icon-image-outline"></i>
                                                    <span></span>
                                                </div>
                                            </el-image>
                                        </div>

                                        <div class="material_detail_wrap" v-if="subItem.data.type == 'video'">
                                            <div class="material_video_icon">
                                                <i class="el-icon-video-camera-solid"></i>
                                                <span class="ml5">{{
                                                    formatDuration(subItem.data.video.duration)
                                                }}</span>
                                            </div>
                                            <video
                                                style="width: 100%;height: 100%;"
                                                :src="subItem.data.video.url"
                                            ></video>
                                        </div>
                                    </div>

                                    <div class="material_intro_wrap">
                                        <div class="material_title">
                                            <div class="material_name" :title="subItem.name">{{ subItem.name }}</div>
                                        </div>
                                        <div class="material_intro">
                                            <span v-if="subItem.data.type == 'image'">
                                                {{ subItem.data.image.width }} x {{ subItem.data.image.height }}px
                                            </span>
                                            <span v-if="subItem.data.type == 'video'">
                                                {{ subItem.data.video.width }} x {{ subItem.data.video.height }}px
                                            </span>
                                        </div>

                                        <div class="material_audit_tag">
                                            <template v-if="subItem.audit_status == 'unaudited'">
                                                <el-tag type="info" size="small">待审核</el-tag>
                                            </template>
                                            <template v-if="subItem.audit_status == 'pass'">
                                                <el-tag type="success" size="small">已通过</el-tag>
                                            </template>
                                            <template v-if="subItem.audit_status == 'reject'">
                                                <el-tooltip
                                                    class="item"
                                                    effect="dark"
                                                    :content="subItem.reject_reason"
                                                    placement="bottom"
                                                >
                                                    <el-tag type="danger" size="small">
                                                        已驳回<i class="ml5 el-icon-question"></i>
                                                    </el-tag>
                                                </el-tooltip>
                                            </template>
                                        </div>

                                        <div class="material_checkbox" v-show="batchEditing">
                                            <el-checkbox
                                                :true-label="subItem.id"
                                                :false-label="0"
                                                v-model="checkBoxItem[subIdx + eachLineCount * idx]"
                                                @change="handleCheckboxchange"
                                            ></el-checkbox>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </el-col>
                    </el-row>
                </template>
            </div>

            <el-divider></el-divider>

            <div class="pagination_wrap">
                <el-pagination
                    class="pagination"
                    @current-change="handleCurrentChange"
                    :current-page="0"
                    :page-size="pageParams.page_size"
                    layout="total, prev, pager, next"
                    :total="total"
                >
                </el-pagination>
            </div>
        </div>

        <div class="material_upload_notice" v-show="uploadingCount > 0">
            <i class="el-icon-loading mr10"></i>
            <span>正在上传{{ uploadingCount }}个任务</span>
            <span class="material_upload_detail_btn" @click="openUploadDetail">【查看详情】</span>
        </div>

        <material-view
            ref="MaterialView"
            :material="selectedItem"
            @reject="reject"
            @pass="pass"
            @delete="doDelete"
        ></material-view>
        <reject ref="reject"></reject>
        <upload ref="upload" :file-list="uploadFileList"></upload>
    </div>
</template>
<script>
import MaterialView from './view';
import Reject from './reject';
import Upload from './upload';
import { Material } from '@dashboard/api';
import _ from 'lodash';

const SCREEN_SIZE_EACH_COUNT = {
    small: {
        min_width: 0,
        max_width: 1000,
        count: 2,
    },
    mid: {
        min_width: 1000,
        max_width: 1500,
        count: 4,
    },
    large: {
        min_width: 1500,
        max_width: Infinity,
        count: 6,
    },
};

// 函数节流
const THROTTLE = function(fn, delay, mustRunDelay) {
    var timer = null;
    var start;
    return function() {
        var context = this,
            args = arguments,
            curr = +new Date();
        clearTimeout(timer);
        if (!start) {
            start = curr;
        }
        if (curr - start >= mustRunDelay) {
            fn.apply(context, args);
            start = curr;
        } else {
            timer = setTimeout(function() {
                fn.apply(context, args);
            }, delay);
        }
    };
};

export default {
    name: 'Material',
    components: {
        MaterialView,
        Reject,
        Upload,
    },
    data() {
        return {
            eachLineCount: 6,

            loading: false,
            list: [],
            listGroup: [[]],
            selectedItem: {},

            total: 0,
            pageParams: {
                page: 1,
                page_size: 30,
            },

            query: {
                name: '',
                // type: '',
                audit_status: '',
            },
            queryOperatorMap: {
                name: 'like',
            },

            checkBoxItem: [],
            checkAll: false,
            indeterminate: false,
            batchEditing: false,

            uploadFileList: [],

            constData: {
                type: [
                    { label: '全部类型', value: '' },
                    { label: '图片', value: 'image' },
                    { label: '视频', value: 'video' },
                ],
                auditStatus: [
                    { label: '全部审核状态', value: '' },
                    { label: '通过', value: 'pass' },
                    { label: '驳回', value: 'reject' },
                    { label: '待审核', value: 'unaudited' },
                ],
            },
        };
    },
    computed: {
        checkedCount() {
            return (this.checkBoxItem || []).filter((item) => {
                return !!item;
            }).length;
        },
        uploadingCount() {
            return (this.uploadFileList || []).filter((item) => {
                return item.percentage != 100;
            }).length;
        },
    },
    methods: {
        getQueryField() {
            let queryFields = [];
            for (let i of Object.keys(this.query)) {
                if (this.query[i] !== '') {
                    let queryField = {
                        field: i,
                        value: this.query[i],
                    };
                    if (this.queryOperatorMap[i]) {
                        queryField.operation = this.queryOperatorMap[i];
                    }
                    queryFields.push(queryField);
                }
            }
            return queryFields;
        },

        //获取素材列表
        getMaterialList(page) {
            const msgPrefix = '拉取素材列表';
            if(page !== undefined){
                this.pageParams.page = page;
            }
            this.loading = true;

            let filter = this.getQueryField();

            Material.get(Object.assign({}, {filter: encodeURIComponent(JSON.stringify(filter))}, this.pageParams))
                .then(({data}) => {
                    this.loading = false;
                    this.total = data.total || 0;
                    this.list = this.formatMaterialList(data.list || []);
                    this.adaptScreenSize();
                    this.resetSelected();
                    this.resetCheckAll();
                })
                .catch(({msg}) => {
                    this.loading = false;
                    this.$message.error(`${msgPrefix}失败 - ${msg}`);
                    this.list = [];
                    this.adaptScreenSize();
                    this.resetSelected();
                    this.resetCheckAll();
                });
        },


        //处理一下数据，防止后台返回数据结果不规范，出现前端报错
        formatMaterialList(list) {
            for (let i = 0; i < list.length; i++) {
                let obj = list[i];
                !obj.data && (obj.data = {});
                !obj.data.image && (obj.data.image = {});
                !obj.data.video && (obj.data.video = {});
            }
            return list;
        },

        //转换时长格式
        formatDuration(duration) {
            duration = duration || 0;
            let s = new Date(duration * 1000).toUTCString();
            s = s.match(/\d{2}:\d{2}:\d{2}/gi);
            s = (s && s[0]) || '00:00:00';
            if (parseInt(duration / 60 / 60) > 0) {
                return s;
            } else {
                return s.slice(3);
            }
        },

        //自适应屏幕尺寸并分行
        adaptScreenSize() {
            let width = document.body.clientWidth;
            let count = 6;
            for (let i of Object.keys(SCREEN_SIZE_EACH_COUNT)) {
                let obj = SCREEN_SIZE_EACH_COUNT[i];
                if (width >= obj.min_width && width < obj.max_width) {
                    count = obj.count;
                    break;
                }
            }
            this.eachLineCount = count;
            let newGroup = [];
            for (let i = 0; i < this.list.length; i++) {
                let lineIdx = parseInt(i / count);
                if (!newGroup[lineIdx]) {
                    newGroup[lineIdx] = [];
                }
                newGroup[lineIdx].push(this.list[i]);
            }
            this.listGroup = newGroup;
        },

        //换页操作
        handleCurrentChange(v) {
            this.getMaterialList(v);
        },

        //选中全部
        handleCheckAllChange(value) {
            if (value) {
                this.checkBoxItem = (this.list || []).map((item) => {
                    return item.id;
                });
            } else {
                for (let i = 0; i < this.checkBoxItem.length; i++) {
                    this.$set(this.checkBoxItem, i, 0);
                }
            }
            this.indeterminate = false;
        },

        handleCheckboxchange() {
            this.resetCheckAll();
        },
        resetSelected() {
            for (let i = 0; i < this.checkBoxItem.length; i++) {
                this.$set(this.checkBoxItem, i, 0);
            }
        },
        resetCheckAll() {
            var count = (this.checkBoxItem || []).filter((item) => {
                return !!item;
            }).length;
            this.indeterminate = count > 0 && count < this.list.length;
            this.checkAll = this.list.length === count;
        },
        //打开预览框
        openPreview(item) {
            this.selectedItem = _.cloneDeep(item);
            this.$refs['MaterialView'].open();
        },

        //检测是否有选中素材
        checkCheckedItem() {
            if (this.checkedCount <= 0) {
                this.$message({
                    message: '请选择素材',
                    type: 'warning',
                });
                return false;
            } else {
                return true;
            }
        },

        //获取选中的素材
        getCheckedItem() {
            return this.checkBoxItem
                .filter((item) => {
                    return !!item;
                })
                .map((item) => {
                    return {
                        id: item,
                    };
                });
        },

        //打开批量操作框
        batchEdit() {
            this.batchEditing = true;
        },
        cancelEdit() {
            this.batchEditing = false;
        },

        //批量驳回
        batchReject() {
            if (!this.checkCheckedItem()) {
                return;
            }
            let list = this.getCheckedItem();
            this.reject(list);
        },

        reject(list) {
            this.$refs['reject'].open().then((rejectReason) => {
                list = list.map((item) => {
                    return {
                        material_id: item.id,
                        audit_status: 'reject',
                        reject_reason: rejectReason,
                    };
                });
                this.loading = true;
                Material.batchAudit(list)
                    .then(() => {
                        this.loading = false;
                        this.$message({
                            message: '素材驳回成功',
                            type: 'success',
                        });
                        this.getMaterialList();
                    })
                    .catch((msg) => {
                        this.loading = false;
                        this.$message.error(`${msg}`);
                    });
            });
        },

        //批量通过
        batchPass() {
            if (!this.checkCheckedItem()) {
                return;
            }
            let list = this.getCheckedItem();
            this.pass(list);
        },
        pass(list) {
            this.$confirm('确认通过选中的素材吗?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning',
            }).then(() => {
                list = list.map((item) => {
                    return {
                        material_id: item.id,
                        audit_status: 'pass',
                    };
                });
                this.loading = true;
                Material.batchAudit(list)
                    .then(() => {
                        this.loading = false;
                        this.$message({
                            message: '素材通过成功',
                            type: 'success',
                        });
                        this.getMaterialList();
                    })
                    .catch((msg) => {
                        this.loading = false;
                        this.$message.error(`${msg}`);
                    });
            });
        },

        //批量删除
        batchDelete() {
            if (!this.checkCheckedItem()) {
                return;
            }
            let list = this.getCheckedItem();
            this.doDelete(list);
        },
        doDelete(list, callback) {
            list = list.map((item) => {
                return item.id;
            });
            this.$confirm('确认删除选中的素材吗?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning',
            }).then(() => {
                this.loading = true;
                Material.batchDelete({ id: list })
                    .then(() => {
                        this.loading = false;
                        this.$message({
                            message: '删除素材成功',
                            type: 'success',
                        });
                        callback && callback();
                        this.getMaterialList();
                    })
                    .catch((msg) => {
                        this.loading = false;
                        this.$message.error(`${msg}`);
                    });
            });
        },

        //上传
        beforeUpload(file) {
            this.uploadFileList.unshift(file);
        },
        replaceFile(file, percentage) {
            for (let i = 0; i < this.uploadFileList.length; i++) {
                if (this.uploadFileList[i].uid == file.uid) {
                    file.percentage = percentage;
                    this.$set(this.uploadFileList, i, file);
                    break;
                }
            }
        },
        checkUploadAllCompleted() {
            let allUploaded =
                (this.uploadFileList || []).findIndex((item) => {
                    return item.percentage != 100;
                }) === -1;
            if (allUploaded) {
                this.getMaterialList(1);
            }
        },
        uploadSuccess(resp, file, fileList) {
            this.replaceFile(file, 50);
            if (!resp.code) {
                this.addMaterial(file, resp.data);
            } else {
                this.replaceFile(file, 100);
            }
        },
        uploadError(err, file, fileList) {
            console.error(err);
            (file.response = {
                code: -1,
                msg: '上传失败',
            }),
                this.replaceFile(file, 100);
            this.checkUploadAllCompleted();
        },

        addMaterial(file, data) {
            let reqData = {
                name: file.name.replace(/\.[a-zA-Z0-9]+$/g, ''),
                data: data,
            };
            Material.edit(reqData)
                .then((data) => {
                    file.response = {
                        code: 0,
                        data: data,
                        msg: '',
                    };
                    this.replaceFile(file, 100);
                    this.checkUploadAllCompleted();
                })
                .catch((e) => {
                    file.response = {
                        code: 1,
                        msg: '上传失败',
                    };
                    this.replaceFile(file, 100);
                    this.checkUploadAllCompleted();
                });
        },

        openUploadDetail() {
            this.$refs['upload'].open();
        },
    },
    created() {
        this.adaptScreenSize();
        window.addEventListener('resize', THROTTLE(this.adaptScreenSize, 500, 6000), false);
    },
    mounted() {
        this.getMaterialList();
    },
};
</script>
<style scoped>
.material_page .material_upload_btn {
    float: right;
}
.material_page .material_wrap {
    border: 1px solid #eee;
    border-radius: 5px;
    height: 240px;
    margin: 15px 15px;
}

.material_page .material_wrap_checked {
    border: 1px solid #409eff;
}

.material_page .material_inner_wrap {
    margin: 5px 5px;
    height: 230px;
}

.material_page .material_detail_wrap {
    margin: 5px 5px;
    background: #eee;
    height: 180px;
    border-radius: 5px;
    position: relative;
    cursor: pointer;
}

.material_page .material_video_icon {
    background: #000;
    opacity: 0.9;
    height: 30px;
    line-height: 30px;
    position: absolute;
    left: 0px;
    top: 0px;
    border-radius: 5px;
    color: #fff;
    font-size: 13px;
    padding: 0px 10px;
}

.material_page .material_video_icon i {
    color: #fff;
    font-size: 20px;
}

.material_page .pagination_wrap {
    height: 50px;
}

.material_page .pagination {
    float: right;
}

.material_page .material_notice {
    text-align: center;
    font-size: 14px;
    color: #aaa;
}

.material_page .material_intro_wrap {
    margin: 5px 5px;
    font-size: 12px;
    line-height: 2em;
    position: relative;
}

.material_page .material_intro_wrap .material_title {
    font-size: 13px;
}

.material_page .material_intro_wrap .material_intro {
    font-size: 12px;
    color: #ccc;
}

.material_page .material_intro_wrap .material_name {
    padding-right: 70px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.material_page .material_intro_wrap .material_audit_tag {
    position: absolute;
    right: 0px;
    top: 0px;
}

.material_page .material_intro_wrap .material_checkbox {
    position: absolute;
    right: 0px;
    bottom: 0px;
}

.material_page .material_handle_bar {
    font-size: 14px;
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    padding: 5px 15px;
    margin-bottom: 20px;
}

.material_page .material_handle_bar .material_handle_close {
    font-weight: bold;
    float: right;
    font-size: 18px;
}

.material_page .material_upload_notice {
    position: fixed;
    width: 600px;
    height: 20px;
    padding: 15px 10px;
    border-radius: 3px;
    border: 1px solid #f2f2f2;
    border-bottom: 4px solid #409eff;
    box-shadow: 0 2px 12px 0 #ccc;
    bottom: 20px;
    left: 50%;
    margin-left: -190px;
    font-size: 13px;
}
.material_page .material_upload_notice .material_upload_detail_btn {
    color: #999;
    float: right;
    cursor: pointer;
}
</style>
<style></style>
