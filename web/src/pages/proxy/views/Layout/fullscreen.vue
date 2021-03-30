<template>
    <el-container class="fullscreen-container">
        <el-header class="fullscreen-header">
            <div class="fullscreen-logo">
                <img src="~@/assets/logo.png" alt="" class="fullscreen-logo-img">
                <h2 class="fullscreen-logo-title">
                    User Growth
                </h2>
            </div>
        </el-header>
        <el-main class="fullscreen-main">
            <slot />
        </el-main>
    </el-container>
</template>
<script>
export default {
    name: 'Fullscreen',
};
</script>
<style lang="scss">
.fullscreen-container {
    position: relative;
    height: 100vh;
    padding: 0;
    margin: 0;
    overflow: hidden;
    background-color: #fefefe;
    .fullscreen-header {
        margin-top: 20px;
    }
    .fullscreen-logo {
        min-height: 64px;
        line-height: 64px;
        .fullscreen-logo-img {
            width: 80px;
            height: 80px;
            vertical-align: middle;
        }
        .fullscreen-logo-title {
            display: inline-block;
            font-weight: 600;
            line-height: 60px;
            font-size: 20px;
            padding-left: 10px;
        }
    }
}
</style>
