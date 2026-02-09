package performance

// Match represents a single match between two alliances of teams.
type Match struct {
	RedTeams  []int
	BlueTeams []int

	RedScore  float64
	BlueScore float64

	RedPenalties  float64
	BluePenalties float64
}

// buildMatchMatrices constructs the matrices A and b used for regression based on the matches and teams.
func buildMatchMatrices(matches []Match, teams []int, scoreFunc func(m Match, isRed bool) float64) ([][]float64, []float64) {
	teamIndex := map[int]int{}
	for i, t := range teams {
		teamIndex[t] = i
	}

	var A [][]float64
	var b []float64

	for _, m := range matches {
		rowRed := make([]float64, len(teams))
		rowBlue := make([]float64, len(teams))

		for _, t := range m.RedTeams {
			rowRed[teamIndex[t]] = 1
		}
		for _, t := range m.BlueTeams {
			rowBlue[teamIndex[t]] = 1
		}

		A = append(A, rowRed)
		b = append(b, scoreFunc(m, true))

		A = append(A, rowBlue)
		b = append(b, scoreFunc(m, false))
	}

	return A, b
}
