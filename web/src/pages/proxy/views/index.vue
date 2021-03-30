<template>
    <fullscreen>
        <el-container class="advertiser-container">
            这是Proxy页面，未承载功能，2秒后为您自动转到管理页面......
        </el-container>
    </fullscreen>
</template>
<script>
import Fullscreen from '@proxy/views/Layout/fullscreen';

export default {
    name: 'Index',
    components: {
        Fullscreen,
    },
    mounted() {
        setTimeout(() => {
            window.location.href = process.env.VUE_APP_PAGE_PORTAL;
        }, 2000);
    },
};
</script>
