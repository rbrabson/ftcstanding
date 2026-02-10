package lambda

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// Match represents a single match between two alliances of teams.
type Match struct {
	RedTeams  []int
	BlueTeams []int
}

// BuildConditionMatrix constructs the design matrix A for OPR/DPR calculations.
// Each row represents an alliance in a match, and each column represents a team.
// The matrix has 1s where a team participated on that alliance, and 0s elsewhere.
func BuildConditionMatrix(matches []Match, teams []int) *mat.Dense {
	// Build team index map
	teamIndex := map[int]int{}
	for i, t := range teams {
		teamIndex[t] = i
	}

	numRows := len(matches) * 2 // Two alliances per match
	numCols := len(teams)

	data := make([]float64, numRows*numCols)
	A := mat.NewDense(numRows, numCols, data)

	row := 0
	for _, m := range matches {
		// Red alliance row
		for _, teamID := range m.RedTeams {
			if idx, ok := teamIndex[teamID]; ok {
				A.Set(row, idx, 1.0)
			}
		}
		row++

		// Blue alliance row
		for _, teamID := range m.BlueTeams {
			if idx, ok := teamIndex[teamID]; ok {
				A.Set(row, idx, 1.0)
			}
		}
		row++
	}

	return A
}

// BuildConditionMatrixFromEvent extracts teams and builds the condition matrix from match data.
func BuildConditionMatrixFromEvent(matches []Match) *mat.Dense {
	// Extract unique teams
	teamSet := make(map[int]struct{})
	for _, m := range matches {
		for _, t := range m.RedTeams {
			teamSet[t] = struct{}{}
		}
		for _, t := range m.BlueTeams {
			teamSet[t] = struct{}{}
		}
	}

	// Convert to sorted slice
	teams := make([]int, 0, len(teamSet))
	for t := range teamSet {
		teams = append(teams, t)
	}

	// Sort teams for consistency
	for i := 0; i < len(teams); i++ {
		for j := i + 1; j < len(teams); j++ {
			if teams[i] > teams[j] {
				teams[i], teams[j] = teams[j], teams[i]
			}
		}
	}

	return BuildConditionMatrix(matches, teams)
}

// AnalyzeEventCondition computes condition number and recommended lambda for an event.
// Returns the condition matrix, its condition number, and the recommended lambda value.
func AnalyzeEventCondition(matches []Match) (a *mat.Dense, condNum float64, lambda float64) {
	a = BuildConditionMatrixFromEvent(matches)

	// Compute condition number of A^T * A
	var ata mat.Dense
	ata.Mul(a.T(), a)
	condNum = ConditionNumber(&ata)

	// Calculate recommended lambda
	matchCount := len(matches)
	lambda = AutoTuneLambda(a, matchCount)

	return a, condNum, lambda
}

func GetLambda(matchCount int) float64 {
	lambda := BaseLambda(matchCount)

	return lambda
}

func BaseLambda(matchCount int) float64 {
	lambda := 0.5 / math.Sqrt(float64(matchCount))

	if lambda < 0.001 {
		return 0.001
	}
	if lambda > 0.3 {
		return 0.3
	}
	return lambda
}

func ConditionNumber(m mat.Matrix) float64 {
	var svd mat.SVD
	ok := svd.Factorize(m, mat.SVDThin)
	if !ok {
		panic("SVD failed")
	}

	values := svd.Values(nil)
	if len(values) == 0 {
		panic("no singular values")
	}

	max := values[0]
	min := values[len(values)-1]

	if min == 0 {
		return math.Inf(1)
	}
	return max / min
}

func RidgeMatrix(a *mat.Dense, lambda float64) *mat.Dense {
	var ata mat.Dense
	ata.Mul(a.T(), a)

	r, _ := ata.Dims()
	for i := 0; i < r; i++ {
		ata.Set(i, i, ata.At(i, i)+lambda)
	}

	return &ata
}

func AutoTuneLambda(a *mat.Dense, matchCount int) float64 {
	const (
		targetCond = 1e7
		maxLambda  = 10.0
	)

	lambda := BaseLambda(matchCount)

	for i := 0; i < 10; i++ {
		M := RidgeMatrix(a, lambda)
		cond := ConditionNumber(M)

		if cond <= targetCond {
			return lambda
		}

		// Increase λ exponentially
		lambda *= 2
		if lambda > maxLambda {
			return maxLambda
		}
	}

	return lambda
}
