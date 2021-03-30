<template>
    <div>
        <el-button type="primary" @click="addGroup">
            新增实验
        </el-button>

        <experiment-group-list />
    </div>
</template>
<script>
import ExperimentGroupList from './list';

export default {
    name: 'ExperimentGroup',
    components: {
        ExperimentGroupList,
    },
    methods: {
        addGroup() {
            this.$router.push({
                name: 'ExperimentGroupEdit',
                params: {
                    id: 0,
                },
            });
        },
    },
};
</script>
