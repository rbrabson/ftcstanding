package main

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/roybrabson/ftcstanding/ftcmath"
)

type Match struct {
	RedTeams  []int
	BlueTeams []int

	RedScore  float64
	BlueScore float64

	RedPenalties  float64
	BluePenalties float64
}

func buildMatrices(
	matches []Match,
	teams []int,
	scoreFunc func(m Match, isRed bool) float64,
) ([][]float64, []float64) {

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

func CalculateNpOPR(matches []Match, teams []int) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.RedPenalties
		}
		return m.BlueScore - m.BluePenalties
	})

	x := ftcmath.SolveLeastSquares(A, b)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func CalculateCCWM(matches []Match, teams []int) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return (m.RedScore - m.BlueScore)
		}
		return (m.BlueScore - m.RedScore)
	})

	x := ftcmath.SolveLeastSquares(A, b)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

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

func CalculateOPRWithReqularization(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore
		}
		return m.BlueScore
	})

	x := ftcmath.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func CalculateNpOPRWithRegularization(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.RedPenalties
		}
		return m.BlueScore - m.BluePenalties
	})

	x := ftcmath.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func CalculateCCWMWithRegularization(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.BlueScore
		}
		return m.BlueScore - m.RedScore
	})

	x := ftcmath.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func CalculateDPR(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			// Red alliance row, opponent is Blue
			return m.BlueScore
		}
		// Blue alliance row, opponent is Red
		return m.RedScore
	})

	x := ftcmath.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func CalculateNpDPR(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.BlueScore - m.BluePenalties
		}
		return m.RedScore - m.RedPenalties
	})

	x := ftcmath.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func CalculateFTCScoutDPR(matches []Match, teams []int, lambda float64, usePenalties bool) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			if usePenalties {
				return m.BlueScore - m.BluePenalties // npDPR
			}
			return m.BlueScore // DPR
		} else {
			if usePenalties {
				return m.RedScore - m.RedPenalties
			}
			return m.RedScore
		}
	})

	x := ftcmath.SolveLeastSquaresRegularized(A, b, lambda)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func LoadMatchesCSV(filename string) ([]Match, []int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	teamSet := map[int]struct{}{}
	matches := []Match{}

	for _, row := range records[1:] { // skip header
		red1, _ := strconv.Atoi(row[0])
		red2, _ := strconv.Atoi(row[1])
		blue1, _ := strconv.Atoi(row[2])
		blue2, _ := strconv.Atoi(row[3])
		redScore, _ := strconv.ParseFloat(row[4], 64)
		blueScore, _ := strconv.ParseFloat(row[5], 64)
		redPen, _ := strconv.ParseFloat(row[6], 64)
		bluePen, _ := strconv.ParseFloat(row[7], 64)

		matches = append(matches, Match{
			RedTeams:      []int{red1, red2},
			BlueTeams:     []int{blue1, blue2},
			RedScore:      redScore,
			BlueScore:     blueScore,
			RedPenalties:  redPen,
			BluePenalties: bluePen,
		})

		teamSet[red1] = struct{}{}
		teamSet[red2] = struct{}{}
		teamSet[blue1] = struct{}{}
		teamSet[blue2] = struct{}{}
	}

	teams := []int{}
	for t := range teamSet {
		teams = append(teams, t)
	}

	return matches, teams, nil
}

func CalculateOPR(matches []Match, teams []int) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore
		}
		return m.BlueScore
	})

	x := ftcmath.SolveLeastSquares(A, b)

	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func CalculateOPRRegularized(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore
		}
		return m.BlueScore
	})
	x := ftcmath.SolveLeastSquaresRegularized(A, b, lambda)
	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func CalculateNpOPRRegularized(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.RedPenalties
		}
		return m.BlueScore - m.BluePenalties
	})
	x := ftcmath.SolveLeastSquaresRegularized(A, b, lambda)
	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}

func CalculateCCWMRegularized(matches []Match, teams []int, lambda float64) map[int]float64 {
	A, b := buildMatrices(matches, teams, func(m Match, isRed bool) float64 {
		if isRed {
			return m.RedScore - m.BlueScore
		}
		return m.BlueScore - m.RedScore
	})
	x := ftcmath.SolveLeastSquaresRegularized(A, b, lambda)
	out := map[int]float64{}
	for i, t := range teams {
		out[t] = x[i]
	}
	return out
}
