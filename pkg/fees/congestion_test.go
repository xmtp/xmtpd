package fees

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXPECTED_CONGESTION_FOR_125 = int64(26)

func TestCalculateAverageCongestion(t *testing.T) {
	rates := [5]int64{125, 100, 100, 100, 100}
	targetRate := int64(100)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, EXPECTED_CONGESTION_FOR_125, congestion)
}

func TestCurrentMinuteCongestion(t *testing.T) {
	rates := [5]int64{125, 0, 0, 0, 0}
	targetRate := int64(100)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, EXPECTED_CONGESTION_FOR_125, congestion)
}

func TestPreviousMinuteCongestion(t *testing.T) {
	rates := [5]int64{0, 125, 0, 0, 0}
	targetRate := int64(100)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, EXPECTED_CONGESTION_FOR_125, congestion)
}

func TestFourMinuteAverageCongestion(t *testing.T) {
	rates := [5]int64{0, 0, 400, 0, 0}
	targetRate := int64(100)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, int64(19), congestion)
}

func TestTargetRateZero(t *testing.T) {
	rates := [5]int64{100, 100, 100, 100, 100}
	targetRate := int64(0)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, int64(0), congestion)
}

func TestWeightedAverage(t *testing.T) {
	rates := [5]int64{100, 100, 100, 100, 200}

	average := calculateFourMinuteAverage(rates)
	assert.Equal(t, int64(110), average)
}
