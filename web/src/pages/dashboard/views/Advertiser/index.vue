<template>
    <div>
        <advertiser-create />

        <advertiser-list @edit="editDialog" />

        <advertiser-edit :ref="advertiserEditRef" :default-value="defaultValue" />
    </div>
</template>
<script>
import AdvertiserList from './list';
import AdvertiserCreate from './create';
import AdvertiserEdit from './edit';

export default {
    name: 'Advertiser',
    components: {
        AdvertiserList,
        AdvertiserCreate,
        AdvertiserEdit,
    },
    data() {
        return {
            defaultValue: {},
            advertiserEditRef: 'advertiserEdit',
        };
    },
    methods: {
        editDialog(data) {
            this.$refs[this.advertiserEditRef].openDialog();
            this.defaultValue = data;
        },
    },
};
</script>
