package fees

import "math"

const (
	// At this point, always return congestion as 100%
	MAX_CONGESTION_RATIO = 1.5
)

/*
*
The algorithm for calculating congestion is as follows:

1) Compute the average number of messages per minute, excluding the current minute which may be incomplete
2) Take the greater of: the current minute's rate, the previous minute's rate, or the preceding 4 minute average
3) Compute the ratio of the max rate to the target rate
4) If the ratio is less than or equal to 1, return 0
5) If the ratio is greater than MAX_CONGESTION_RATIO, return 100
6) Otherwise, return the congestion as an exponential function of the ratio
*
*/
func CalculateCongestion(last5Minutes [5]int32, targetRatePerMinute int32) int32 {
	// If the target rate is 0, return 0
	if targetRatePerMinute == 0 {
		return 0
	}

	currentMinute := last5Minutes[0]
	prevMinute := last5Minutes[1]

	// 1) Compute the average congestion, excluding the current minute which may be incomplete
	fourMinuteAverage := calculateFourMinuteAverage(last5Minutes)

	// 2) Determine if congestion must be 0:
	//    Congestion is 0 only if both the current and previous minute
	//    are not above target AND the 4-minute average is not above target.
	if currentMinute <= targetRatePerMinute && prevMinute <= targetRatePerMinute &&
		fourMinuteAverage <= targetRatePerMinute {
		return 0
	}

	// 3) Compute the ratio = max(current, previous, 4-min avg) / T
	ratio := math.Max(
		float64(currentMinute),
		math.Max(float64(prevMinute), float64(fourMinuteAverage)),
	) / float64(
		targetRatePerMinute,
	)

	// 4) If ratio <= 1, effectively no excess above T
	//    (In practice, if we reached here, ratio should be >= 1.)
	if ratio <= 1 {
		return 0
	}

	// 5) If ratio >= MAX_CONGESTION_RATIO => congestion = 100
	if ratio >= MAX_CONGESTION_RATIO {
		return 100
	}

	// 6) Exponential mapping for 1 < ratio < MAX_CONGESTION_RATIO
	//    - ratio=1   => 0
	//    - ratio=MAX_CONGESTION_RATIO => 100
	//    - Grows exponentially in between
	k := 4.0 // Adjust k to tune how fast congestion rises
	numerator := math.Exp(k*(ratio-1.0)) - 1.0
	denominator := math.Exp(k*(MAX_CONGESTION_RATIO-1.0)) - 1.0
	congestion := 100.0 * (numerator / denominator)

	return int32(congestion)
}

func calculateFourMinuteAverage(last5Minutes [5]int32) int32 {
	// Apply weights that decrease with recency (index 1 is most recent, index 4 is oldest)
	// Weight distribution: 40%, 30%, 20%, 10%
	weightedSum := last5Minutes[1]*4 + last5Minutes[2]*3 + last5Minutes[3]*2 + last5Minutes[4]*1
	totalWeight := int32(10) // sum of weights: 4+3+2+1
	return weightedSum / totalWeight
}
