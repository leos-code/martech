package loader

import "testing"

func TestParseHDFSPath(t *testing.T) {
	t.Log(ParseHDFSPath("hdfs://ss-cdg-3-v2/stage/outface/sng/attribution/enckeys/enckeys"))
}
