package api

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/brahma-adshonor/gohook"
	"github.com/stretchr/testify/assert"
)

func Test_generateAuthorization(t *testing.T) {
	rtaID := "100011502"
	token := "dfcd79dbbd08a0203dbde9a9cc59"
	var now int64 = 1497196800

	assert.Equal(t, "32f0a07d6884f179113ff144c7c988be", generateAuthorization(rtaID, token, now))
}

func TestFillRTAExpRequestHeader(t *testing.T) {
	err := gohook.Hook(time.Now, func() time.Time {
		return time.Unix(1497196800, 0)
	}, nil)
	assert.NoError(t, err)

	rtaID := "100011502"
	token := "dfcd79dbbd08a0203dbde9a9cc59"

	header := http.Header{}
	FillRTAExpRequestHeader(header, rtaID, token)
	assert.EqualValues(t, rtaID, header.Get("RtaId"))
	assert.EqualValues(t, strconv.FormatInt(1497196800, 10), header.Get("Time"))
	assert.EqualValues(t, "application/json", header.Get("Content-Type"))
	assert.EqualValues(t, "32f0a07d6884f179113ff144c7c988be", header.Get("Authorization"))
}
