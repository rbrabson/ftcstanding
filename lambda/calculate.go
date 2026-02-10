package lambda

import (
	"github.com/rbrabson/ftcstanding/performance"
)

// GetLambda computes the recommended lambda value for regularization based on the match data.
func GetLambda(matches []performance.Match) float64 {
	lambda := baseLambda(len(matches))
	return lambda
}

// BaseLambda computes the base lambda value based on the number of matches.
func baseLambda(matchCount int) float64 {
	switch {
	case matchCount < 20:
		return 0.1
	case matchCount <= 60:
		return 0.01
	default:
		return 0.001
	}
}
