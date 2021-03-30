package novelty

import (
	"github.com/tencentad/martech/api/proto/novelty"
)

func AcceptClick(data *novelty.Novelty, clickTime int64) {
	day := clickTime / 86400

	click := data.Click
	click.AccumulateClickCount += 1
	if click.LastClickDay == day {
		click.DayClickCount += 1
	} else {
		click.DayClickCount = 1
		click.LastClickDay = day
	}
	click.LastClickTimestamp = clickTime
}
