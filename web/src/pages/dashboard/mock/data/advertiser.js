export const AdvertiserList = {
    code: 0,
    data: [
        {
            id: 1,
            created_at: '2021-02-25T11:01:29.69013+08:00',
            updated_at: '2021-02-25T11:01:29.69013+08:00',
            platform: 'ams',
            name: 'AMS账号 A',
            object_id: 45,
            object: {
                id: 45,
                created_at: '2021-02-25T11:09:58.253458+08:00',
                updated_at: '2021-02-25T11:09:58.253458+08:00',
                name: '账号主体 A',
                type: 'advertiser',
                object: 'object',
                tenant_id: 10,
                parent_id: 36,
            },
            tenant_id: 10,
        },
        {
            id: 1,
            created_at: '2021-02-25T11:01:29.69013+08:00',
            updated_at: '2021-02-25T11:01:29.69013+08:00',
            platform: 'oceanengine',
            name: '巨量引擎账号 B',
            object_id: 46,
            object: {
                id: 46,
                created_at: '2021-02-25T11:10:09.747808+08:00',
                updated_at: '2021-03-01T20:41:11.710199+08:00',
                name: '账号主体 B',
                type: 'advertiser',
                object: 'object_33',
                tenant_id: 10,
                parent_id: 36,
            },
            tenant_id: 10,
        },
    ],
    msg: 'get success',
};

export const AuthorityList = {
    code: 0,
    data: [
        {
            platform: 'ams',
            url: 'https://developers.e.qq.com/oauth/authorize?client_id=1111186039&redirect_uri=http%3A%2F%2Fdev.ug.com%2Fams%3Fuid%3D1__OBJECT_ID__',
        },
        {
            platform: 'oceanengine',
            url: 'https://www.baidu.com/__OBJECT_ID__',
        },
    ],
    msg: 'get success',
}
