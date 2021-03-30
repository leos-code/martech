export const targetingSchema = {
    code: 0,
    data: {
        "version": 1,
        "fields":
            [
                {
                    "name": "age",
                    "display_name": "年龄",
                    "type": "integer",
                    "range": {
                        "min": 0,
                        "max": 200
                    }
                },
                {
                    "name": "last_active_time",
                    "display_name": "上次活跃时间",
                    "type": "enum",
                    "enum": [{"value": "1周内"}, {"value": "1周以外，1月以内"}, {"value": "1月以外"}]
                }, {"name": "self_active_rate", "display_name": "自启概率", "type": "integer"},
                {
                    "name": "interest",
                    "display_name": "兴趣爱好",
                    "type": "enum",
                    "enum": [
                        {
                            "value": "游戏",
                            "children": [
                                {
                                    "value": "王者荣耀",
                                    "children": []
                                },
                                {
                                    "value": "和平精英"
                                },
                                {
                                    "value": "阴阳师",
                                    // children: [
                                    //     {"value": "阴阳师1"},
                                    //     {"value": "阴阳师2"}
                                    // ]
                                }
                            ]
                        },
                        {
                            "value": "音乐"
                        }
                    ]
                },
                {
                    "name": "download_song",
                    "display_name": "下载歌曲",
                    "type": "string"
                }
            ]
    },
    msg: 'success'
}
