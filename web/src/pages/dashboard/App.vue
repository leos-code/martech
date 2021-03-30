<template>
    <div id="app">
        <ug-layout />
    </div>
</template>

<script>
import '@/style/base.scss';
import '@/style/main.scss';

import UgLayout from '@dashboard/views/Layout/layout';

export default {
    name: 'Dashboard',
    components: {
        UgLayout,
    },
};
</script>
