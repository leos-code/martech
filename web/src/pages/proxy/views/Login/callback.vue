<template>
    <div />
</template>
<script>
import Utils from '@/utils';

export default {
    name: 'LoginCallback',
    mounted() {
        if (window.top && Utils.isFunction(window.top.loginCallback)) {
            window.top.loginCallback(this.$route.query);
        }
    },
};
</script>
