<template>
    <div>您已成功完成注册, 3秒后自动关闭当前页面......</div>
</template>
<script>
export default {
    name: 'AdvertiserProxy',
    mounted() {
        setTimeout(() => {
            window.opener = null;
            window.open('', '_self');
            window.close();
        }, 3000);
    },
};
</script>
