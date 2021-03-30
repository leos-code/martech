package novelty

import (
	"testing"

	"github.com/tencentad/martech/api/proto/novelty"
	"github.com/stretchr/testify/assert"
)

func TestAcceptClick(t *testing.T) {
	{
		data := &novelty.Novelty{
			Click: &novelty.Click{
				DayClickCount:        1,
				LastClickDay:         100,
				AccumulateClickCount: 10,
				LastClickTimestamp:   8640001,
			},
		}

		AcceptClick(data, 8640002)
		assert.EqualValues(t, 2, data.Click.DayClickCount)
		assert.EqualValues(t, 11, data.Click.AccumulateClickCount)
		assert.EqualValues(t, 100, data.Click.LastClickDay)
		assert.EqualValues(t, 8640002, data.Click.LastClickTimestamp)
	}

	{
		data := &novelty.Novelty{
			Click: &novelty.Click{
				DayClickCount:        1,
				LastClickDay:         100,
				AccumulateClickCount: 10,
				LastClickTimestamp:   8640001,
			},
		}

		AcceptClick(data, 86400000)
		assert.EqualValues(t, 1, data.Click.DayClickCount)
		assert.EqualValues(t, 1000, data.Click.LastClickDay)
		assert.EqualValues(t, 11, data.Click.AccumulateClickCount)
		assert.EqualValues(t, 86400000, data.Click.LastClickTimestamp)
	}

}
