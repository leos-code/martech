import { superAdminTenantId, localUserId, localRoleId } from './meta';

export const userInfo = {
    code: 0,
    data: {
        user: {
            id: localUserId,
            created_at: '2021-02-03T10:33:48.375218+08:00',
            updated_at: '2021-02-03T10:33:48.375218+08:00',
            phone_number: '18506663368',
            email: '562276119@qq.com',
        },
        login_user: {
            open_id: '56672',
            login_type: 'rio',
            created_at: '2021-02-03T10:33:48.375652+08:00',
            updated_at: '2021-02-03T10:33:48.375652+08:00',
            user_id: localUserId,
            avatar: 'http://r.hrc.oa.com/photo/150/dickzheng.png',
            nick_name: 'dickzheng',
        },
        current_tenant: {
            id: superAdminTenantId,
            name: 'God管理台',
        },
        tenant: [
            {
                id: superAdminTenantId,
                name: 'God管理台',
            },
        ],
    },
    msg: 'get success',
};

export const userList = {
    code: 0,
    data: [
        {
            id: 1,
            created_at: '2021-02-04T10:29:23.342884+08:00',
            updated_at: '2021-02-04T10:29:23.342884+08:00',
            phone_number: '18506663368',
            email: '562276119@qq.com',
            login_user: [
                {
                    open_id: '56672',
                    login_type: 'rio',
                    created_at: '2021-02-04T10:29:23.346777+08:00',
                    updated_at: '2021-02-04T10:29:23.346777+08:00',
                    user_id: 1,
                    avatar: 'http://r.hrc.oa.com/photo/150/dickzheng.png',
                    nick_name: 'dickzheng',
                },
            ],
            role: [
                {
                    id: 1,
                    name: 'Role1',
                    parent_id: 3,
                    user_list: [1, 2],
                },
                {
                    id: 3,
                    name: 'Role3',
                    parent_id: 0,
                    user_list: [4],
                },
            ],
        },
        {
            id: 2,
            created_at: '2021-02-04T10:29:23.342884+08:00',
            updated_at: '2021-02-04T10:29:23.342884+08:00',
            phone_number: '18888888888',
            email: '562276119@qq.com',
            login_user: [
                {
                    open_id: '56672',
                    login_type: 'rio',
                    created_at: '2021-02-04T10:29:23.346777+08:00',
                    updated_at: '2021-02-04T10:29:23.346777+08:00',
                    user_id: 2,
                    avatar: 'http://r.hrc.oa.com/photo/150/linceyou.png',
                    nick_name: 'linceyou',
                },
            ],
        },
        {
            id: 3,
            created_at: '2021-02-04T10:29:23.342884+08:00',
            updated_at: '2021-02-04T10:29:23.342884+08:00',
            phone_number: '18888888888',
            email: '562276119@qq.com',
            login_user: [
                {
                    open_id: '56672',
                    login_type: 'rio',
                    created_at: '2021-02-04T10:29:23.346777+08:00',
                    updated_at: '2021-02-04T10:29:23.346777+08:00',
                    user_id: 3,
                    avatar: 'http://r.hrc.oa.com/photo/150/slyao.png',
                    nick_name: 'slyao',
                },
            ],
        },
        {
            id: 4,
            created_at: '2021-02-04T10:29:23.342884+08:00',
            updated_at: '2021-02-04T10:29:23.342884+08:00',
            phone_number: '18888888888',
            email: '562276119@qq.com',
            login_user: [
                {
                    open_id: '56672',
                    login_type: 'rio',
                    created_at: '2021-02-04T10:29:23.346777+08:00',
                    updated_at: '2021-02-04T10:29:23.346777+08:00',
                    user_id: 4,
                    avatar: 'http://r.hrc.oa.com/photo/150/alnwang.png',
                    nick_name: 'alnwang',
                },
            ],
        },
        {
            id: 5,
            created_at: '2021-02-04T10:29:23.342884+08:00',
            updated_at: '2021-02-04T10:29:23.342884+08:00',
            phone_number: '18888888888',
            email: '562276119@qq.com',
            login_user: [
                {
                    open_id: '56672',
                    login_type: 'rio',
                    created_at: '2021-02-04T10:29:23.346777+08:00',
                    updated_at: '2021-02-04T10:29:23.346777+08:00',
                    user_id: 5,
                    avatar: 'http://r.hrc.oa.com/photo/150/zelfaliu.png',
                    nick_name: 'zelfaliu',
                },
            ],
        },
    ],
    msg: 'get success',
};
