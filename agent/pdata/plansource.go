package pdata

import "github.com/chaosblade-io/chaosblade/data"

type PSource struct {
	data.Source
}

var pSource PSource

func init() {
	source := data.GetSource().(*data.Source)
	pSource = PSource{*source}
	pSource.CheckAndInitExperimentPlanTable()
}

func GetSource() PSource {
	return pSource
}

const tableExistsDQL = `SELECT count(*) AS c
	FROM sqlite_master 
	WHERE type = "table"
	AND name = ?
`