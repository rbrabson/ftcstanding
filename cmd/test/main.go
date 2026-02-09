package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rbrabson/ftcstanding/performance"
)

func GetLambda(numTeams int) float64 {
	return 1.0 / float64(numTeams)
}

func LoadMatchesCSV(filename string) ([]performance.Match, []int, error) {
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
	matches := []performance.Match{}

	for _, row := range records[1:] { // skip header
		red1, _ := strconv.Atoi(row[0])
		red2, _ := strconv.Atoi(row[1])
		blue1, _ := strconv.Atoi(row[2])
		blue2, _ := strconv.Atoi(row[3])
		redScore, _ := strconv.ParseFloat(row[4], 64)
		blueScore, _ := strconv.ParseFloat(row[5], 64)
		redPen, _ := strconv.ParseFloat(row[6], 64)
		bluePen, _ := strconv.ParseFloat(row[7], 64)

		matches = append(matches, performance.Match{
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

func main() {
	matches, teams, err := LoadMatchesCSV("matches.csv")
	if err != nil {
		log.Fatal(err)
	}

	lambda := GetLambda(len(teams)) // FTCScout-style regularization

	opr := performance.CalculateOPR(matches, teams)
	npopr := performance.CalculateNpOPR(matches, teams)
	ccwm := performance.CalculateCCWM(matches, teams)
	dpr := performance.CalculateDPR(matches, teams, lambda)
	npdpr := performance.CalculateNpDPR(matches, teams, lambda)

	fmt.Println("Team | OPR   | npOPR | CCWM  | DPR  | npDPR | npAVG")
	fmt.Println("----------------------------------------------------")
	for _, t := range teams {
		npavg := performance.CalculateNpAVG(matches, t)
		fmt.Printf("%4d | %5.2f | %5.2f | %5.2f | %5.2f | %5.2f | %5.2f\n",
			t, opr[t], npopr[t], ccwm[t], dpr[t], npdpr[t], npavg)
	}

	opr = performance.CalculateNpOPRWithRegularaization(matches, teams, lambda)
	npopr = performance.CalculateNpOPRWithRegularization(matches, teams, lambda)
	ccwm = performance.CalculateCCWMWithRegularization(matches, teams, lambda)
	dpr = performance.CalculateDPRWithRegularization(matches, teams, lambda)
	npdpr = performance.CalculateNpDPRWithRegularization(matches, teams, lambda)

	fmt.Println()
	fmt.Println("Team | OPR   | npOPR | CCWM  | DPR   | npDPR | npAVG")
	fmt.Println("----------------------------------------------------")
	for _, t := range teams {
		npavg := performance.CalculateNpAVG(matches, t)
		fmt.Printf("%4d | %5.2f | %5.2f | %5.2f | %5.2f | %5.2f | %5.2f\n",
			t, opr[t], npopr[t], ccwm[t], dpr[t], npdpr[t], npavg)
	}
}
