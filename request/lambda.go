package request

import (
	"github.com/rbrabson/ftcstanding/performance"
)

// getLambda computes the recommended lambda value for regularization based on the match data.
func getLambda(matches []performance.Match) float64 {
	matcheCount := len(matches)
	switch {
	case matcheCount < 20:
		return 0.1
	case matcheCount <= 60:
		return 0.01
	default:
		return 0.001
	}
}
