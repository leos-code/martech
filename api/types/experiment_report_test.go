package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddExperimentReportMetric(t *testing.T) {
	m := AddExperimentReportMetric(nil, 1)
	assert.EqualValues(t, 1, m.Value)
	m1 := AddExperimentReportMetric(m, 2)
	assert.EqualValues(t, 3, m1.Value)
}

func TestExperimentReportRecord_CalculateDerivedMetric(t *testing.T) {
	{
		r := &ExperimentReportRecord{
			Cost:             &ExperimentReportMetric{Value: 100},
			Exposure:         &ExperimentReportMetric{Value: 100000},
			Click:            &ExperimentReportMetric{Value: 1000},
			Conversion:       &ExperimentReportMetric{Value: 80},
			ConversionSecond: &ExperimentReportMetric{Value: 50},
		}

		r.CalculateDerivedMetric()

		assert.EqualValues(t, 1, r.Cpm.Value)
		assert.EqualValues(t, 0.01, r.Ctr.Value)
		assert.EqualValues(t, 0.1, r.Cpc.Value)
		assert.EqualValues(t, 0.08, r.Cvr.Value)
		assert.EqualValues(t, 0.05, r.CvrSecond.Value)
	}

	{
		r := &ExperimentReportRecord{
			Cost:             &ExperimentReportMetric{Value: 100},
			Exposure:         &ExperimentReportMetric{Value: 0},
			Click:            &ExperimentReportMetric{Value: 0},
			Conversion:       &ExperimentReportMetric{Value: 80},
			ConversionSecond: &ExperimentReportMetric{Value: 50},
		}

		r.CalculateDerivedMetric()

		assert.EqualValues(t, 0, r.Cpm.Value)
		assert.EqualValues(t, 0, r.Ctr.Value)
		assert.EqualValues(t, 0, r.Cpc.Value)
		assert.EqualValues(t, 0, r.Cvr.Value)
		assert.EqualValues(t, 0, r.CvrSecond.Value)
	}

}

func TestCalculateExperimentReportRecordDelta(t *testing.T) {
	{
		CalculateExperimentReportRecordDelta(nil)
	}
	{
		r1 := &ExperimentReportRecord{
			Cost:             &ExperimentReportMetric{Value: 100},
			Exposure:         &ExperimentReportMetric{Value: 100000},
			Click:            &ExperimentReportMetric{Value: 1000},
			Conversion:       &ExperimentReportMetric{Value: 80},
			ConversionSecond: &ExperimentReportMetric{Value: 50},
		}

		r1.CalculateDerivedMetric()

		r2 := &ExperimentReportRecord{
			Cost:             &ExperimentReportMetric{Value: 100},
			Exposure:         &ExperimentReportMetric{Value: 100000},
			Click:            &ExperimentReportMetric{Value: 500},
			Conversion:       &ExperimentReportMetric{Value: 80},
			ConversionSecond: &ExperimentReportMetric{Value: 50},
		}
		r2.CalculateDerivedMetric()

		CalculateExperimentReportRecordDelta([]*ExperimentReportRecord{r1, r2})

		assert.Nil(t, r1.Cvr.Delta)
		assert.EqualValues(t, 0, *r2.Cost.Delta)
		assert.EqualValues(t, 0, *r2.Exposure.Delta)
		assert.EqualValues(t, 0, *r2.Conversion.Delta)
		assert.EqualValues(t, 0, *r2.ConversionSecond.Delta)
		assert.EqualValues(t, -0.5, *r2.Click.Delta)
		assert.EqualValues(t, 0, *r2.Cpm.Delta)
		assert.EqualValues(t, -0.5, *r2.Ctr.Delta)
		assert.EqualValues(t, 1, *r2.Cpc.Delta)
		assert.EqualValues(t, 1, *r2.Cvr.Delta)
		assert.EqualValues(t, 1, *r2.CvrSecond.Delta)
	}
}
