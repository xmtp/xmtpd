package fees

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedCongestionFor125 = int32(26)

func TestCalculateAverageCongestion(t *testing.T) {
	rates := [5]int32{125, 100, 100, 100, 100}
	targetRate := int32(100)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, expectedCongestionFor125, congestion)
}

func TestCurrentMinuteCongestion(t *testing.T) {
	rates := [5]int32{125, 0, 0, 0, 0}
	targetRate := int32(100)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, expectedCongestionFor125, congestion)
}

func TestPreviousMinuteCongestion(t *testing.T) {
	rates := [5]int32{0, 125, 0, 0, 0}
	targetRate := int32(100)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, expectedCongestionFor125, congestion)
}

func TestFourMinuteAverageCongestion(t *testing.T) {
	rates := [5]int32{0, 0, 400, 0, 0}
	targetRate := int32(100)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, int32(19), congestion)
}

func TestTargetRateZero(t *testing.T) {
	rates := [5]int32{100, 100, 100, 100, 100}
	targetRate := int32(0)

	congestion := CalculateCongestion(rates, targetRate)
	assert.Equal(t, int32(0), congestion)
}

func TestWeightedAverage(t *testing.T) {
	rates := [5]int32{100, 100, 100, 100, 200}

	average := calculateFourMinuteAverage(rates)
	assert.Equal(t, int32(110), average)
}
