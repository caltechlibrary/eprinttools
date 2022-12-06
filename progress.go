package eprinttools

import (
	"fmt"
	"time"
)

func progress(t0 time.Time, i int, tot int) string {
	if i == 0 {
		return "0.00 ETA Unknown"
	}
	// percent completed
	percent := (float64(i) / float64(tot)) * 100.0
	if i == 0 {
	}
	// running time
	rt := time.Now().Sub(t0)
	// estimated time remaining
	eta := time.Duration((float64(rt) / float64(i) * float64(tot)) - float64(rt))
	return fmt.Sprintf("%.2f%% ETA %v", percent, eta.Truncate(time.Second))
}
