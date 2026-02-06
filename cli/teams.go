package cli

import (
	"bytes"

	"github.com/fatih/color"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rodaine/table"
)

// RenderTeams renders a list of teams as a table.
func RenderTeams(teams []*database.Team) string {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "Country", "Region")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, team := range teams {
		tbl.AddRow(team.TeamID, team.Name, team.Country, team.HomeRegion)
	}

	buffer := &bytes.Buffer{}
	tbl.WithWriter(buffer).Print()
	return buffer.String()
}
