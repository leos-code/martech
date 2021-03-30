package deviceid

// ID 设备ID
type ID struct {
	IDFA         string `json:"idfa"`
	IDFAMD5      string `json:"idfa_md5"`
	IMEI         string `json:"imei"`
	IMEIMD5      string `json:"imei_md5"`
	OAID         string `json:"oaid"`
	OAIDMD5      string `json:"oaid_md5"`
	AndroidID    string `json:"android_id"`
	AndroidIDMD5 string `json:"android_id_md5"`
	MAC          string `json:"mac"`
	MACMD5       string `json:"mac_md5"`
	QAID         string `json:"qaid"`
}
