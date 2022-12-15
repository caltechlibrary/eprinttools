package eprinttools

import (
	"fmt"
	"time"
)

func calcModValue(tot int) int {
	switch {
		case tot < 100:
			return 10
		case tot < 1000:
			return 100
		case tot < 10000:
			return 500
		case tot < 25000:
			return 1000
		case tot < 50000:
			return 1500
		case tot < 100000:
			return 2000
		default:
			return 2500
	}
}

func progress(t0 time.Time, i int, tot int) string {
	if i == 0 {
		return "calc. Estimated Time Remaining"
	}
	// percent completed
	percent := (float64(i) / float64(tot)) * 100.0
	if i == 0 {
	}
	// running time
	rt := time.Now().Sub(t0)
	// estimated time remaining
	eta := time.Duration((float64(rt) / float64(i) * float64(tot)) - float64(rt))
	return fmt.Sprintf("%.2f%% ETR %v", percent, eta.Truncate(time.Second))
}
