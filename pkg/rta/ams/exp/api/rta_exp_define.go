package api

var (
	// RtaExpRespCodeMap 错误码对应的错误信息
	RtaExpRespCodeMap = map[int]string{
		1601: "请携带签名信息",
		1602: "签名校验失败",
		1603: "授权Token已过期",
		1604: "头信息缺少RTA ID",
		1605: "头信息缺少时间戳",
		1606: "请求已过期，请检查时间戳",
		1607: "RTA ID不合法",
		1608: "Token未设置，请设置后再请求",
		1:    "通用错误",
	}

	RtaExpSiteSetMap = map[int]string{
		15: "移动联盟",
		21: "微信",
		25: "移动内部站点",
		27: "腾讯新闻流量",
		28: "腾讯视频流量",
	}
)
