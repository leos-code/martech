export const targetingList = {
    "code": 0, "data": {
        "total": 1,
        "list": [
            {
                "id": 6,
                "created_at": "2021-02-22T08:45:52.53Z",
                "updated_at": "2021-02-23T02:42:36.824Z",
                "name": "aaa",
                "targeting_info": [{
                    "name": "age",
                    "not": false,
                    "values": {"type": "range", "range": [{"start": 11, "end": 22}, {"start": 33, "end": 44}]}
                }, {
                    "name": "interest",
                    "not": false,
                    "values": {"type": "string", "string": ["王者荣耀", "和平精英", "阴阳师"]}
                }, {"name": "download_song", "not": false, "values": {"type": "string", "string": ["abcd"]}}],
                "bind_strategy": [{
                    "name": "aa",
                    "platform": "ams",
                    "strategy": {"strategy_id": ["bb"], "advertiser_id": ["dd"], "campaign_id": ["ee"]}
                }, {
                    "name": "bb",
                    "platform": "ams",
                    "strategy": {"strategy_id": ["hh", "jj"], "advertiser_id": ["kk"], "campaign_id": ["ll"]}
                }],
                "freq_control": {
                    rules:[{
                        "platform": "ams",
                        "for": "click",
                        "limit": 11,
                        "accumulate_limit": 22,
                        "interval": 33
                    }]
                }
            }
        ]
    }, "msg": "get success"
}
export const targetingEdit = {"code":0,"msg":"success"}
