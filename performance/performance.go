package performance

import "github.com/rbrabson/ftcstanding/matrix"

// CalculateCCWM calculates the Calculated Contribution to Winning Margin (CCWM) for each team.
func CalculateCCWM(matches []Match, teams []int) map[int]float64 {
	A, b := buildMatchMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return (m.RedScore - m.BlueScore)
		}
		return (m.BlueScore - m.RedScore)
	})

	x := matrix.SolveLeastSquares(A, b)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

// CalculateDPR calculates a team's Defensive Power Rating (DPR) using ridge regression with regularization.
func CalculateDPR(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatchMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.BlueScore
		}
		return m.RedScore
	})

	x := matrix.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

// CalculateNpAVG calculates the average non-penalized score for a given team across all matches.
func CalculateNpAVG(matches []Match, team int) float64 {
	var total float64
	var count float64

	for _, m := range matches {
		for _, t := range m.RedTeams {
			if t == team {
				total += m.RedScore - m.RedPenalties
				count++
			}
		}
		for _, t := range m.BlueTeams {
			if t == team {
				total += m.BlueScore - m.BluePenalties
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}
	return total / count
}

// CalculateNpDPR calculates a team's non-penalized Defensive Power Rating (DPR) using ridge regression with regularization.
func CalculateNpDPR(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatchMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.BlueScore - m.BluePenalties
		}
		return m.RedScore - m.RedPenalties
	})

	x := matrix.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

// CalculateNpOPR calculates a team's non-penalized Offensive Power Rating (OPR).
func CalculateNpOPR(matches []Match, teams []int) map[int]float64 {
	A, b := buildMatchMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.RedPenalties
		}
		return m.BlueScore - m.BluePenalties
	})

	x := matrix.SolveLeastSquares(A, b)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

// CalculateOPR calculates a team's Offensive Power Rating (OPR).
func CalculateOPR(matches []Match, teams []int) map[int]float64 {
	A, b := buildMatchMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore
		}
		return m.BlueScore
	})

	x := matrix.SolveLeastSquares(A, b)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}
