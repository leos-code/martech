package common

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/godruid"
	"time"
)

func parseTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05-0700")
}

func BuildIntervals(t1, t2 time.Time) []string {
	intervals := fmt.Sprintf("%v/%v", parseTime(t1), parseTime(t2))
	return []string{intervals}
}

func BuildDimensions(dimensions []interface{}) []godruid.DimSpec {
	ds := make([]godruid.DimSpec, 0)
	linq.From(dimensions).ToSlice(&ds)
	return ds
}
