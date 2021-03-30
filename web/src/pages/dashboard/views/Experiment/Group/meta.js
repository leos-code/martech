export const EXPERIMENT_STAGE_STATUS = {
    Draft: {
        key: 'Draft',
        name: '草稿',
    },
    Running: {
        key: 'Running',
        name: '运行中',
    },
    Stop: {
        key: 'Stop',
        name: '停止',
    },
};

export const RULES = {
    group: {
        name: [
            {
                required: true,
                message: '请填写实验组名称',
                trigger: 'change',
            },
        ],
        rta_account_id: [
            {
                required: true,
                message: '请选择RTA账户',
                trigger: 'change',
            },
        ],
        draft: {
            end_time: [
                {
                    required: true,
                    message: '请选择结束时间',
                    trigger: 'change',
                },
            ],
        },
    },
    expItem: {
        name: [
            {
                required: true,
                message: '请输入实验名称',
                trigger: 'change',
            },
        ],
        rta_exp: [
            {
                type: 'array',
                required: true,
                message: '请绑定rta exp',
                trigger: 'change',
            },
        ],
    },
};
