<template>
    <el-dialog title="素材预览" :visible.sync="visible" width="1000px" :append-to-body="true" :modal-append-to-body="true"
               :close-on-click-modal="false" @closed="onClose" @opened="onOpen">
        <div>
            <template v-if="!!materialShow">
                <el-row :gutter="30">
                    <el-col :span="16">
                        <template v-if="materialShow.data.type == 'image'">
                            <el-card class="material_left_card" ref="material_wrap">
                                <el-image
                                    class="material_image"
                                    :src="materialShow.data.image.url"
                                    fit="scale-down"
                                    :preview-src-list="[materialShow.data.image.url]"
                                    :z-index="10000"
                                >
                                    <div slot="error" class="image-slot">
                                        <i class="el-icon-picture-outline"></i>
                                    </div>
                                </el-image>
                            </el-card>
                        </template>
                        <template v-if="materialShow.data.type == 'video'">
                            <el-card class="material_left_card" ref="material_wrap">
                                <video poster=""
                                       ref="video"
                                       class="material_video"
                                       :src="materialShow.data.video.url"
                                       controls="controls"
                                       :style="videoStyle"
                                       autoplay="autoplay"
                                >
                                </video>
                            </el-card>
                        </template>
                    </el-col>
                    <el-col :span="8">
                        <el-card class="material_right_card">
                            <p class="material_intro_line">
                                <span class="material_label">标题：</span>
                                <span class="material_value">{{materialShow.name}}</span>
                            </p>
                            <template v-if="materialShow.data.type == 'image'">
                                <p class="material_intro_line">
                                    <span class="material_label">尺寸：</span>
                                    <span>{{materialShow.data.image.width}} x {{materialShow.data.image.height}} px</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">大小：</span>
                                    <span>{{formatSize(materialShow.data.image.size)}}</span>
                                </p>
                            </template>
                            <template v-if="materialShow.data.type == 'video'">
                                <p class="material_intro_line">
                                    <span class="material_label">分辨率：</span>
                                    <span>{{materialShow.data.video.width}} x {{materialShow.data.video.height}} px</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">大小：</span>
                                    <span>{{formatSize(materialShow.data.video.size)}}</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">视频时长：</span>
                                    <span>{{materialShow.data.video.duration}} 秒</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">格式：</span>
                                    <span>{{materialShow.data.video.container_format}}</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">视频编码格式：</span>
                                    <span>{{materialShow.data.video.codec_format}}</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">码率：</span>
                                    <span>{{materialShow.data.video.bitrate}} kbit/s</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">帧率：</span>
                                    <span>{{materialShow.data.video.frame_rate}}</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">音频格式：</span>
                                    <span>{{materialShow.data.video.audio_codec}}</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">音频码率：</span>
                                    <span>{{materialShow.data.video.audio_bitrate}} kbit/s</span>
                                </p>
                                <p class="material_intro_line">
                                    <span class="material_label">音频采样率：</span>
                                    <span>{{materialShow.data.video.audio_sample_rate}} kHZ</span>
                                </p>
                            </template>
                        </el-card>
                    </el-col>
                </el-row>
            </template>
        </div>
        <span slot="footer" class="dialog-footer">
            <el-button @click="visible = false">关 闭</el-button>
            <el-button type="danger" @click="doDelete" v-permission="'MaterialEdit'">删除此素材</el-button>
            <el-button type="danger" @click="reject" v-permission="'MaterialAudit'">驳 回</el-button>
            <el-button type="success" @click="pass" v-permission="'MaterialAudit'">通 过</el-button>
        </span>

    </el-dialog>
</template>

<script>
    import _ from 'lodash'

    export default {
        name: "material_view",
        props: {
            material: Object
        },
        data() {
            return {
                visible: false,

                materialShow: undefined,

                videoStyle: {
                    width: 'auto',
                    height: 'auto',
                }
            }
        },
        methods: {
            onClose() {
                this.materialShow = undefined;
            },
            beforeOpen(){

            },
            onOpen() {
                this.materialShow = this.material || {};
                this.$nextTick(() => {
                    let video = this.$refs['video'];
                    if (video) {
                        this.videoStyle = {
                            width: 'auto',
                            height: 'auto',
                        };
                        if (video.readyState == 4) {
                            this.autoScaleDownVideo(video);
                        } else {
                            let handler = ()=>{
                                video.removeEventListener('canplaythrough', handler)
                                this.autoScaleDownVideo(video);
                            };
                            video.addEventListener('canplaythrough',handler)
                        }
                    }
                })
            },
            open() {
                this.visible = true;
            },

            formatSize(size) {
                size = parseInt(size / 1024);
                if (size < 1024) {
                    return _.round(size, 2) + ' KB';
                } else {
                    size = parseInt(size) / 1024;
                    return _.round(size, 2) + ' MB';
                }
            },
            formatDuration(duration) {
                duration = duration || 0;
                let s = new Date(duration * 1000).toUTCString()
                s = s.match(/\d{2}:\d{2}:\d{2}/gi);
                s = s && s[0] || '00:00:00'
                console.log(duration);
                if (parseInt(duration / 60 / 60) > 0) {
                    return s
                }else{
                    return s.slice(3);
                }
            },

            //自适应视频尺寸
            autoScaleDownVideo(video) {
                let width = video.videoWidth;
                let height = video.videoHeight;
                let wrap = this.$refs['material_wrap'];
                let wrapWidth = wrap.$el.clientWidth - 40;
                let wrapHeight = wrap.$el.clientHeight - 40;
                var widthScale = width / wrapWidth;
                var heightScale = height / wrapHeight;
                if (widthScale > 1 || heightScale > 1) {
                    if (widthScale > heightScale) {
                        width = wrapWidth;
                        height = height / widthScale;
                    } else {
                        height = wrapHeight;
                        width = width / heightScale;
                    }
                }
                this.videoStyle = {width: parseInt(width)+'px', height: parseInt(height) + 'px'}
            },

            reject(){
                this.$emit('reject', [{id: this.material.id}]);
            },
            pass(){
                this.$emit('pass', [{id: this.material.id}]);
            },
            doDelete(){
                this.$emit('delete', [{id: this.material.id}], ()=>{
                    this.visible = false;
                });
            }
        },
        mounted() {

        }
    }
</script>

<style scoped>
    .material_left_card {
        width: 100%;
        height: 440px;
        background: #EEE;
        box-shadow: none;
        display: flex;
        display: -webkit-flex;
        align-items:center;
        justify-content:center;
    }
    .material_left_card .material_image {
        width: 100%;
        height: 400px;
    }

    .material_left_card .material_video {
        outline: 0 none !important;
    }

    .material_right_card {
        height: 440px;
        font-size: 12px;
        line-height: 1.8em;
        overflow: auto;
    }

    .material_right_card .material_intro_line{
        padding-left: 90px;
        position: relative;
    }

    .material_right_card .material_label{
        width: 90px;
        color: #AAA;
        position: absolute;
        top: 0;
        left: 0;

    }
    .material_right_card .material_value{
        display: inline-block;
    }
</style>
