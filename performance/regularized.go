package performance

import "github.com/rbrabson/ftcstanding/matrix"

// CalculateCCWMWithRegularization calculates the Calculated Contribution to Winning Margin (CCWM) for each team using ridge regression with regularization.
func CalculateCCWMWithRegularization(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatchMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.BlueScore
		}
		return m.BlueScore - m.RedScore
	})

	x := matrix.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

// CalculateDPRWithRegularization calculates a team's Defensive Power Rating (DPR) using ridge regression with regularization.
func CalculateDPRWithRegularization(matches []Match, teams []int, lambda float64) map[int]float64 {
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

// CalculateNpDPRWithRegularization calculates a team's non-penalized Defensive Power Rating (DPR) using ridge regression with regularization.
func CalculateNpDPRWithRegularization(matches []Match, teams []int, lambda float64) map[int]float64 {
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

// CalculateNpOPRWithRegularization calculates a team's non-penalized Offensive Power Rating (OPR) using ridge regression with regularization.
func CalculateNpOPRWithRegularaization(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatchMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.RedPenalties
		}
		return m.BlueScore - m.BluePenalties
	})
	x := matrix.SolveLeastSquaresRegularized(A, b, lambda)
	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

// CalculateNpOPRWithRegularization calculates a team's non-penalized Offensive Power Rating (OPR) using ridge regression with regularization.
func CalculateNpOPRWithRegularization(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatchMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.RedPenalties
		}
		return m.BlueScore - m.BluePenalties
	})

	x := matrix.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

// CalculateOPRWithRegularization calculates a team's Offensive Power Rating (OPR) using ridge regression with regularization.
func CalculateOPRWithRegularization(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatchMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore
		}
		return m.BlueScore
	})

	x := matrix.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}
