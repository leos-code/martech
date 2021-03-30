<template>
    <fullscreen>
        <el-container class="advertiser-container">
            您已成功完成注册, 3秒后自动关闭当前页面......
        </el-container>
    </fullscreen>
</template>
<script>
import Fullscreen from '@proxy/views/Layout/fullscreen';

export default {
    name: 'Advertiser',
    components: {
        Fullscreen,
    },
    mounted() {
        setTimeout(() => {
            window.opener = null;
            window.open('', '_self');
            window.close();
        }, 3000);
    },
};
</script>
