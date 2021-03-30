<template>
    <div />
</template>
<script>
export default {
    name: 'Login',
    mounted() {
        window.location.href = process.env.VUE_APP_PAGE_PORTAL;
    },
};
</script>
