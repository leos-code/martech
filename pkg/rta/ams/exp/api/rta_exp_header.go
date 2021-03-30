package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// RtaExpRequestBase 请求媒体侧RTA实验系统基础参数
type RtaExpRequestBase struct {
	RtaID string
	Token string
}

// FillRTAExpRequestHeader 填充RTA实验请求的http header
func FillRTAExpRequestHeader(header http.Header, rtaID string, token string) {
	header.Set("RtaId", rtaID)
	now := time.Now().Unix()
	header.Set("Time", strconv.FormatInt(now, 10))
	header.Set("Authorization", generateAuthorization(rtaID, token, now))
	header.Set("Content-Type", "application/json")
}

func generateAuthorization(rtaID string, token string, time int64) string {
	input := rtaID + token + strconv.FormatInt(time, 10)
	fmt.Println(input)
	m := md5.New()
	m.Write([]byte(input))
	return hex.EncodeToString(m.Sum(nil))
}
