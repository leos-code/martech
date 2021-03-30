export const materialList = {
    "code": 0, "data": {
        "total": 1,
        "list": [{
            "id": 6,
            "created_at": "2021-02-22T08:45:52.53Z",
            "updated_at": "2021-02-23T02:42:36.824Z",
            "name": "147*147",
            "data": {
                "type": "image",
                "image": {
                    "url": "http://ug-1258344696.cos.ap-guangzhou.myqcloud.com/material/20210222084552_147%2A147.png",
                    "size": 15135,
                    "ext": ".png",
                    "width": 147,
                    "height": 147
                }
            },
            "audit": [{
                "id": 1,
                "created_at": "2021-02-23T02:42:22.098Z",
                "updated_at": "2021-02-23T02:42:22.098Z",
                "material_id": 6,
                "material": null,
                "user_id": 11,
                "user": null,
                "audit_status": "pass"
            }, {
                "id": 3,
                "created_at": "2021-02-23T02:42:36.823Z",
                "updated_at": "2021-02-23T02:42:36.823Z",
                "material_id": 6,
                "material": null,
                "user_id": 11,
                "user": null,
                "audit_status": "reject",
                "reject_reason": "test"
            }],
            "audit_status": "reject",
            "reject_reason": "test"
        }, {
            "id": 7,
            "created_at": "2021-02-22T08:46:04.966Z",
            "updated_at": "2021-02-23T02:42:36.827Z",
            "name": "660*346",
            "data": {
                "type": "image",
                "image": {
                    "url": "http://ug-1258344696.cos.ap-guangzhou.myqcloud.com/material/20210222084604_660%2A346.png",
                    "size": 39646,
                    "ext": ".png",
                    "width": 1320,
                    "height": 696
                }
            },
            "audit": [{
                "id": 2,
                "created_at": "2021-02-23T02:42:22.101Z",
                "updated_at": "2021-02-23T02:42:22.101Z",
                "material_id": 7,
                "material": null,
                "user_id": 11,
                "user": null,
                "audit_status": "pass"
            }, {
                "id": 4,
                "created_at": "2021-02-23T02:42:36.826Z",
                "updated_at": "2021-02-23T02:42:36.826Z",
                "material_id": 7,
                "material": null,
                "user_id": 11,
                "user": null,
                "audit_status": "reject",
                "reject_reason": "test"
            }],
            "audit_status": "reject",
            "reject_reason": "test"
        }]
    }, "msg": "get success"
};

export const materialUpload = {
    code: 0,
    data: {
        "code": 0,
        "data": {
            "type": "image",
            "image": {
                "url": "http://ug-1258344696.cos.ap-guangzhou.myqcloud.com/material/20210222070944_160%2A160.jpg",
                "size": 5169,
                "ext": ".jpg",
                "width": 160,
                "height": 160
            }
        },
        "msg": "upload success"
    },
    msg: 'success',
}
// export const materialUpload = {"code":-1,"msg":"exec: \"ffprobe\": executable file not found in $PATH"}

export const materialEdit = {
    code: 0,
    data: {},
    msg: 'success',
}

export const materialSubmitAudit = {
    code: 0,
    data: {},
    msg: 'success',
}
