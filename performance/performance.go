package performance

import "github.com/rbrabson/ftcstanding/matrix"

// Calculator calculates various performance metrics for teams based on match data.
type Calculator struct {
	Matches []Match
	Teams   []int
	Lambda  float64
}

// CalculateCCWM calculates the Calculated Contribution to Winning Margin (CCWM) for each team.
func (p *Calculator) CalculateCCWM() map[int]float64 {
	A, b := buildMatchMatrices(p.Matches, p.Teams, func(m Match, isRed bool) float64 {
		if isRed {
			return (m.RedScore - m.BlueScore)
		}
		return (m.BlueScore - m.RedScore)
	})

	var x []float64
	if p.Lambda == 0 {
		x = matrix.SolveLeastSquares(A, b)
	} else {
		x = matrix.SolveLeastSquaresRegularized(A, b, p.Lambda)
	}

	out := map[int]float64{}
	for i, t := range p.Teams {
		out[t] = x[i]
	}
	return out
}

// CalculateDPR calculates the Defensive Power Rating (DPR) for each team.
func (p *Calculator) CalculateDPR() map[int]float64 {
	A, b := buildMatchMatrices(p.Matches, p.Teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.BlueScore
		}
		return m.RedScore
	})
	var x []float64
	if p.Lambda == 0 {
		x = matrix.SolveLeastSquares(A, b)
	} else {
		x = matrix.SolveLeastSquaresRegularized(A, b, p.Lambda)
	}

	out := map[int]float64{}
	for i, t := range p.Teams {
		out[t] = x[i]
	}
	return out
}

// CalculateNpAVG calculates the non-penalized average score for a given team.
func (p *Calculator) CalculateNpAVG(matches []Match, team int) float64 {
	var total float64
	var count float64

	for _, m := range p.Matches {
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

// CalculateNpDPR calculates the non-penalized Defensive Power Rating (DPR) for each team.
func (p *Calculator) CalculateNpDPR() map[int]float64 {
	A, b := buildMatchMatrices(p.Matches, p.Teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.BlueScore - m.BluePenalties
		}
		return m.RedScore - m.RedPenalties
	})

	var x []float64
	if p.Lambda == 0 {
		x = matrix.SolveLeastSquares(A, b)
	} else {
		x = matrix.SolveLeastSquaresRegularized(A, b, p.Lambda)
	}

	out := map[int]float64{}
	for i, t := range p.Teams {
		out[t] = x[i]
	}
	return out

}

// CalculateNpOPR calculates the non-penalized Offensive Power Rating (OPR) for each team.
func (p *Calculator) CalculateNpOPR() map[int]float64 {
	A, b := buildMatchMatrices(p.Matches, p.Teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.RedPenalties
		}
		return m.BlueScore - m.BluePenalties
	})

	var x []float64
	if p.Lambda == 0 {
		x = matrix.SolveLeastSquares(A, b)
	} else {
		x = matrix.SolveLeastSquaresRegularized(A, b, p.Lambda)
	}

	out := map[int]float64{}
	for i, t := range p.Teams {
		out[t] = x[i]
	}
	return out
}

// CalculateOPR calculates the Offensive Power Rating (OPR) for each team.
func (p *Calculator) CalculateOPR() map[int]float64 {
	A, b := buildMatchMatrices(p.Matches, p.Teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore
		}
		return m.BlueScore
	})

	var x []float64
	if p.Lambda == 0 {
		x = matrix.SolveLeastSquares(A, b)
	} else {
		x = matrix.SolveLeastSquaresRegularized(A, b, p.Lambda)
	}

	out := map[int]float64{}
	for i, t := range p.Teams {
		out[t] = x[i]
	}
	return out
}
