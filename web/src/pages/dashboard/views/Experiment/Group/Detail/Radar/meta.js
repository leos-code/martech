/**
 * 时间粒度
 * dateFormat: YYYY-MM-DD HH:mm:ss
 */
export const Granularity = [
    {
        key: 'min_5',
        name: '5分钟',
        range: 3,
        dateFmt: 'HH:MM',
    },
    {
        key: 'hour',
        name: '小时',
        range: 7,
        dateFmt: 'mm月dd日 HH:00',
    },
    {
        key: 'day',
        name: '天',
        range: 90,
        dateFmt: 'mm月dd日',
    },
];
export const Dimentions = [
    {
        key: 'time',
        name: '时间',
    },
    {
        key: 'experiment_item.name',
        name: '实验组',
    },
];
export const Fields = [
    {
        key: 'cost',
        name: '消耗',
        tips: '总消耗，包括合约广告',
        values: [
            {
                key: 'value',
                formatter: [
                    {
                        type: 'round',
                        config: {
                            precision: 3,
                        },
                    },
                    {
                        type: 'thousandSeparator',
                    },
                ],
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
    {
        key: 'exposure',
        name: '曝光',
        values: [
            {
                key: 'value',
                formatter: [
                    {
                        type: 'round',
                        config: {
                            precision: 3,
                        },
                    },
                    {
                        type: 'thousandSeparator',
                    },
                ],
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
    {
        key: 'click',
        name: '点击',
        values: [
            {
                key: 'value',
                formatter: [
                    {
                        type: 'round',
                        config: {
                            precision: 3,
                        },
                    },
                    {
                        type: 'thousandSeparator',
                    },
                ],
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
    {
        key: 'cpm',
        name: 'CPM',
        values: [
            {
                key: 'value',
                formatter: [
                    {
                        type: 'round',
                        config: {
                            precision: 3,
                        },
                    },
                    {
                        type: 'thousandSeparator',
                    },
                ],
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
    {
        key: 'cpc',
        name: 'CPC',
        values: [
            {
                key: 'value',
                formatter: [
                    {
                        type: 'round',
                        config: {
                            precision: 3,
                        },
                    },
                    {
                        type: 'thousandSeparator',
                    },
                ],
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
    {
        key: 'ctr',
        name: 'CTR',
        values: [
            {
                key: 'value',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
    {
        key: 'conversion',
        name: '浅层转化量',
        values: [
            {
                key: 'value',
                formatter: [
                    {
                        type: 'round',
                        config: {
                            precision: 3,
                        },
                    },
                    {
                        type: 'thousandSeparator',
                    },
                ],
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
    {
        key: 'conversion_second',
        name: '深层转化量',
        tips: '深层转化量(第二目标)',
        values: [
            {
                key: 'value',
                formatter: [
                    {
                        type: 'round',
                        config: {
                            precision: 3,
                        },
                    },
                    {
                        type: 'thousandSeparator',
                    },
                ],
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
    {
        key: 'cvr',
        name: '浅层CVR',
        values: [
            {
                key: 'value',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
    {
        key: 'cvr_second',
        name: '深层CVR',
        tips: '深层CVR(第二目标)',
        values: [
            {
                key: 'value',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
            },
            {
                key: 'delta',
                formatter: {
                    type: 'percentage',
                    config: {
                        precision: 3,
                    },
                },
                colorful: true,
            },
        ],
    },
];
