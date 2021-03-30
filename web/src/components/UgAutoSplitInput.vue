<template>
    <span>
        <span v-if="label" class="label">
            {{ label }}
        </span>
        <mixin-select
            class="ug-auto-split-input"
            :placeholder="placeHolder"
            :value="v"
            @change="updateValue"
            v-loading="loading"
            allow-create
            multiple
            filterable
            clearable
            :visible="false"
            popper-class="ug-auto-split-input-option"
            :split-char="splitChar"
            :size="size"
        />
    </span>
</template>
<script>
import _ from 'lodash';
import MixinSelect from './MixinSelect';

export default {
    name: 'UgAutoSplitInput',
    components: {
        MixinSelect,
    },
    props: {
        splitChar: {
            type: Array,
            default: () => {
                return [',', '，', ';', '\n'];
            },
        },
        label: String,
        placeholder: String,
        size: String,
        value: Array,
    },
    watch: {
        label() {
            this.init();
        },
        value(v){
            this.v = v;
        }
    },
    computed: {
        placeHolder() {
            if (this.placeholder) {
                return this.placeholder;
            } else if (this.splitChar && this.splitChar.length > 0) {
                const splitChar = _.cloneDeep(this.splitChar);
                const seperator = splitChar
                    .map((item) => {
                        switch (item) {
                            case '\n':
                                return '回车键';
                            case '\t':
                                return '制表符';
                            default:
                                return item;
                        }
                    })
                    .join(' ');
                return `分隔符为 ${seperator}`;
            } else {
                return '';
            }
        },
    },
    data() {
        return {
            loading: false,
            v: '',
        };
    },
    methods: {
        init() {},
        updateValue(v) {
            this.v = v;
            this.$emit('input', v);
        },
    },
    created() {},
    mounted() {
        this.updateValue(this.value)
    },
};
</script>
<style>
.ug-auto-split-input .el-icon-arrow-up {
    display: none;
}
.ug-auto-split-input-option {
    display: none;
}
</style>
