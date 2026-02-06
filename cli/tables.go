package cli

import (
	"bytes"

	"github.com/fatih/color"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rodaine/table"
)

func RenderTable(db database.DB) string {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "Country", "Region")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	filter := database.TeamFilter{
		HomeRegions: []string{"USNC"},
	}
	for _, team := range db.GetAllTeams(filter) {
		tbl.AddRow(team.TeamID, team.Name, team.Country, team.HomeRegion)
	}

	buffer := &bytes.Buffer{}
	tbl.WithWriter(buffer).Print()
	return buffer.String()
}
