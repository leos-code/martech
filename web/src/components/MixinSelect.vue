<template>
    <div
        class="el-select"
        :class="[selectSize ? 'el-select--' + selectSize : '']"
        @click.stop="toggleMenu"
        v-clickoutside="handleClose"
    >
        <div
            class="el-select__tags"
            v-if="multiple"
            ref="tags"
            :style="{ 'max-width': inputWidth - 32 + 'px', width: '100%' }"
        >
            <span v-if="collapseTags && selected.length">
                <el-tag
                    :closable="!selectDisabled"
                    :size="collapseTagSize"
                    :hit="selected[0].hitState"
                    type="info"
                    @close="deleteTag($event, selected[0])"
                    disable-transitions
                >
                    <span class="el-select__tags-text">{{ selected[0].currentLabel }}</span>
                </el-tag>
                <el-tag
                    v-if="selected.length > 1"
                    :closable="false"
                    :size="collapseTagSize"
                    type="info"
                    disable-transitions
                >
                    <span class="el-select__tags-text">+ {{ selected.length - 1 }}</span>
                </el-tag>
            </span>
            <transition-group @after-leave="resetInputHeight" v-if="!collapseTags">
                <el-tag
                    v-for="item in selected"
                    :key="getValueKey(item)"
                    :closable="!selectDisabled"
                    :size="collapseTagSize"
                    :hit="item.hitState"
                    type="info"
                    @close="deleteTag($event, item)"
                    disable-transitions
                >
                    <span class="el-select__tags-text">
                        {{ item.currentLabel }}
                    </span>
                </el-tag>
            </transition-group>

            <textarea
                class="el-select__input"
                :class="[selectSize ? `is-${selectSize}` : '']"
                :disabled="selectDisabled"
                :autocomplete="autoComplete || autocomplete"
                @focus="handleFocus"
                @blur="softFocus = false"
                @click.stop
                @keyup="managePlaceholder"
                @keydown="resetInputState"
                @keydown.down.prevent="navigateOptions('next')"
                @keydown.up.prevent="navigateOptions('prev')"
                @keydown.enter.prevent="myMultiSelectOption"
                @keydown.esc.stop.prevent="visible = false"
                @keydown.delete="deletePrevTag"
                @compositionstart="handleComposition"
                @compositionupdate="handleComposition"
                @compositionend="handleComposition"
                v-model="query"
                @input="myDebouncedQueryChange"
                v-if="filterable"
                :style="{
                    'flex-grow': '1',
                    resize: 'none',
                    'margin-left': '15px',
                    'margin-top': '12px',
                    font: '400 13.3333px Arial',
                    width: inputLength / (inputWidth - 32) + '%',
                    'max-width': inputWidth - 42 + 'px',
                    height: '24px',
                }"
                ref="input"
            />
        </div>
        <el-input
            ref="reference"
            v-model="selectedLabel"
            type="text"
            :placeholder="currentPlaceholder"
            :name="name"
            :id="id"
            :autocomplete="autoComplete || autocomplete"
            :size="selectSize"
            :disabled="selectDisabled"
            :readonly="readonly"
            :validate-event="false"
            :class="{ 'is-focus': visible }"
            @focus="handleFocus"
            @blur="handleBlur"
            @keyup.native="debouncedOnInputChange"
            @keydown.native.down.stop.prevent="navigateOptions('next')"
            @keydown.native.up.stop.prevent="navigateOptions('prev')"
            @keydown.native.enter.prevent="mySelectOption"
            @keydown.native.esc.stop.prevent="visible = false"
            @keydown.native.tab="visible = false"
            @paste.native="debouncedOnInputChange"
            @mouseenter.native="inputHovering = true"
            @mouseleave.native="inputHovering = false"
        >
            <template slot="prefix" v-if="$slots.prefix">
                <slot name="prefix"></slot>
            </template>
            <template slot="suffix">
                <i v-show="!showClose" :class="['el-select__caret', 'el-input__icon', 'el-icon-' + iconClass]"></i>
                <i
                    v-if="showClose"
                    class="el-select__caret el-input__icon el-icon-circle-close"
                    @click="handleClearClick"
                ></i>
            </template>
        </el-input>
        <transition name="el-zoom-in-top" @before-enter="handleMenuEnter" @after-leave="doDestroy">
            <el-select-menu ref="popper" :append-to-body="popperAppendToBody" v-show="visible && emptyText !== false">
                <el-scrollbar
                    tag="ul"
                    wrap-class="el-select-dropdown__wrap"
                    view-class="el-select-dropdown__list"
                    ref="scrollbar"
                    :class="{ 'is-empty': !allowCreate && query && filteredOptionsCount === 0 }"
                    v-show="options.length > 0 && !loading"
                >
                    <el-option :value="query" created v-if="showNewOption"></el-option>
                    <slot></slot>
                </el-scrollbar>
                <template v-if="emptyText && (!allowCreate || loading || (allowCreate && options.length === 0))">
                    <slot name="empty" v-if="$slots.empty"></slot>
                    <p class="el-select-dropdown__empty" v-else>{{ emptyText }}</p>
                </template>
            </el-select-menu>
        </transition>
    </div>
</template>
<script>
import { Select } from 'element-ui';

export default {
    name: 'MixinSelect',
    props: {
        splitChar: {
            type: Array,
            default: () => {
                return [];
            },
        },
    },
    mixins: [Select],
    methods: {
        myDebouncedQueryChange(event) {
            // 处理输入事件
            // 配置splitChar，则判断是否需要split
            if (this.splitChar.length >= 1) {
                if (event.inputType === 'insertFromPaste') {
                    // 粘贴操作
                    let re = this.splitChar[0];
                    for (let i = 1; i < this.splitChar.length; i++) {
                        re = re.concat('|' + (this.splitChar[i] === '|' ? '\\|' : this.splitChar[i]));
                    }
                    let arrs = this.query.split(new RegExp(`${re}`)).filter((item) => item !== '');
                    const value = (this.value || []).slice();
                    for (let i = 0; i < arrs.length; i++) {
                        const optionIndex = this.getValueIndex(value, arrs[i]);
                        if (optionIndex > -1) {
                            // 若值已存在，则不需要操作
                            // value.splice(optionIndex, 1);
                        } else if (this.multipleLimit <= 0 || value.length < this.multipleLimit) {
                            // 判断是否可以create
                            if (this.allowCreate) {
                                // 可以create直接添加值
                                value.push(arrs[i]);
                            } else {
                                // 不可以create则需要option存在该值
                                let options = this.options || [];
                                let optionsIndex = options.findIndex((item) => item.visible && item.value === arrs[i]);
                                if (optionsIndex > -1) {
                                    value.push(arrs[i]);
                                }
                            }
                        }
                    }
                    this.$emit('input', value);
                    this.emitChange(value);
                } else if (event.inputType === 'insertText') {
                    // 输入操作
                    let splitChar = this.splitChar ? this.splitChar.find((item) => item === event.data) : undefined;
                    if (splitChar) {
                        this.query = this.query.substring(0, this.query.length - splitChar.length);
                        if (this.setFirstHoverIndex() > -1) {
                            this.selectOption(event);
                        }
                    }
                }
            }
            // 原生的处理方法
            this.debouncedQueryChange(event);
        },
        myMultiSelectOption(event) {
            if (this.splitChar.length === 0 || this.splitChar.includes('\n')) {
                if (this.setFirstHoverIndex() > -1) {
                    this.selectOption(event);
                }
            }
        },
        mySelectOption(event) {
            if (this.setFirstHoverIndex() > -1) {
                this.selectOption(event);
            }
        },
        setFirstHoverIndex() {
            //未手动选中任何项目，需处理
            if (this.hoverIndex <= -1) {
                let options = this.options || [];
                let selectedValue = undefined;
                let selectedIdx = -1;
                let visibleOptions = options.filter((item) => item.visible);
                if (!visibleOptions.length) {
                    this.hoverIndex = -1;
                } else {
                    if (this.allowCreate && options[options.length - 1].created === true) {
                        selectedValue = options[options.length - 1].value;
                        selectedIdx = options.length - 1;
                    } else {
                        for (let i = 0; i < this.options.length; i++) {
                            if (this.options[i].visible) {
                                selectedValue = this.options[i].value;
                                selectedIdx = i;
                                break;
                            }
                        }
                    }
                }
                if (selectedIdx > -1) {
                    if (
                        (this.multiple && this.value.includes(selectedValue)) ||
                        (!this.multiple && this.value === selectedValue)
                    ) {
                        this.hoverIndex = -1;
                    } else {
                        this.hoverIndex = selectedIdx;
                    }
                } else {
                    this.hoverIndex = -1;
                }
            }
            return this.hoverIndex;
        },
        getOption(value) {
            let option;
            const isObject = Object.prototype.toString.call(value).toLowerCase() === '[object object]';
            const isNull = Object.prototype.toString.call(value).toLowerCase() === '[object null]';
            const isUndefined = Object.prototype.toString.call(value).toLowerCase() === '[object undefined]';

            for (let i = 0; i <= this.cachedOptions.length - 1; i++) {
                const cachedOption = this.cachedOptions[i];
                const isEqual = isObject
                    ? // eslint-disable-next-line no-undef
                      getValueByPath(cachedOption.value, this.valueKey) === getValueByPath(value, this.valueKey)
                    : cachedOption.value === value;
                if (isEqual) {
                    option = cachedOption;
                    break;
                }
            }
            if (option) return option;
            const label = !isObject && !isNull && !isUndefined ? value : '';
            let newOption = {
                value: value,
                currentLabel: label,
            };
            if (this.multiple) {
                newOption.hitState = false;
            }
            return newOption;
        },
    },
};
</script>
