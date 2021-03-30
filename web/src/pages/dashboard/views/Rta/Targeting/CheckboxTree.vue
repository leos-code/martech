<template>
    <span>
        <span v-if="treeLevel <= 1">
            <template v-if="isValidData">
                <el-checkbox-group v-model="outputList" @change="handleChange">
                    <el-checkbox v-for="item in data" :key="item.value" :label="item.value">{{ item.value }}</el-checkbox>
                </el-checkbox-group>
            </template>
        </span>
        <span v-else :class="treeLevel === 2 ? ['targeting_two_level_tree']: []" style="margin-top: 20px;">
            <template v-if="isValidData">
                <el-tree
                    ref="tree"
                    :data="data"
                    :props="{children: 'children',label: 'value'}"
                    :default-expand-all="true"
                    :expand-on-click-node="false"
                    :indent="treeLevel === 2 ? 48 : 24"
                    class="targeting_tree"
                    show-checkbox
                    node-key="value"
                    @node-click="handleClickNode"
                    @check-change="handleChange"
                />
            </template>
        </span>
    </span>
</template>

<script>
    import _ from 'lodash';
    export default {
        name: "CheckboxGroups",
        props: {
            data: Array,
            value: {
                type: Array,
            }
        },
        watch: {
            value(v){
                //若长度不同或存在不一样的元素，则认为不一致
                if(v && ((v.length !== this.outputList.length) || (_.difference(v, this.outputList).length !== 0))){
                    this.setCheckedKeys();
                }
            }
        },
        data() {
            return {
                checkedList: [],
                outputList: [],
            }
        },
        computed: {
            isValidData() {
                return _.isArray(this.data)
            },
            treeLevel() {
                if (!this.isValidData) {
                    return 1;
                } else {
                    let f = (nodes, level) => {
                        return nodes.reduce((a, b) => {
                            return Math.max(a.children && _.isArray(a.children) && a.children.length > 0 ? f(a.children, level + 1) : level, b.children && _.isArray(b.children) && b.children.length > 0 ? f(b.children, level + 1) : level)
                        });
                    };
                    return f(this.data, 1);
                }
            }
        },
        methods: {
            handleClickNode(node) {
                let checked = this.checkedList.findIndex(item => item === node.value) > -1;
                this.$refs['tree'].setChecked(node.value, !checked, true);
            },
            handleChange() {
                this.$nextTick(() => {
                    if(this.treeLevel > 1){
                        this.checkedList = this.$refs['tree'].getCheckedKeys();
                        this.outputList = this.getCheckedKeys();
                    }
                    this.$emit('input', this.outputList)
                })
            },
            getCheckedKeys(){
                return this.$refs['tree'].getCheckedNodes().filter((item)=>{
                    return !item.children || !item.children.length;
                }).map((item)=>{
                    return item.value
                })
            },
            setCheckedKeys(){
                if(this.treeLevel > 1){
                    this.$refs['tree'].setCheckedKeys(this.value || []);
                }
                this.outputList = this.value || [];
                this.$nextTick(()=>{
                    this.handleChange();
                })
            }
        },
        mounted() {
            this.setCheckedKeys();
        }
    }
</script>

<style>
    .targeting_tree{
        margin-top: 5px;
    }
    .targeting_two_level_tree .el-tree-node__children {
        white-space: normal;
        margin-top: 10px;
    }

    .targeting_two_level_tree .el-tree-node__children .el-tree-node {
        display: inline-block;
        white-space: normal;
    }

    .targeting_two_level_tree .el-tree-node__children .el-tree-node .el-tree-node__content:hover,
    .targeting_two_level_tree .el-tree-node__children .el-tree-node.is-current .el-tree-node__content:hover {
        background: none !important;
    }
</style>
