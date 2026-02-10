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
// It only includes teams that actually participate in the provided matches to reduce matrix size.
// Returns: A matrix, b vector, and list of active teams (in order corresponding to matrix columns).
func buildMatchMatrices(matches []Match, teams []int, scoreFunc func(m Match, isRed bool) float64) ([][]float64, []float64, []int) {
	// First, identify which teams actually participate in these matches
	participatingTeams := make(map[int]struct{})
	for _, m := range matches {
		for _, t := range m.RedTeams {
			participatingTeams[t] = struct{}{}
		}
		for _, t := range m.BlueTeams {
			participatingTeams[t] = struct{}{}
		}
	}

	// Create a list of participating teams in sorted order (to match the teams list order)
	var activeTeams []int
	for _, t := range teams {
		if _, ok := participatingTeams[t]; ok {
			activeTeams = append(activeTeams, t)
		}
	}

	// Build index map for active teams only
	teamIndex := make(map[int]int)
	for i, t := range activeTeams {
		teamIndex[t] = i
	}

	var a [][]float64
	var b []float64

	for _, m := range matches {
		rowRed := make([]float64, len(activeTeams))
		rowBlue := make([]float64, len(activeTeams))

		for _, t := range m.RedTeams {
			if idx, ok := teamIndex[t]; ok {
				rowRed[idx] = 1
			}
		}
		for _, t := range m.BlueTeams {
			if idx, ok := teamIndex[t]; ok {
				rowBlue[idx] = 1
			}
		}

		a = append(a, rowRed)
		b = append(b, scoreFunc(m, true))

		a = append(a, rowBlue)
		b = append(b, scoreFunc(m, false))
	}

	return a, b, activeTeams
}
