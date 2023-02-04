package image

import (
	"bytes"
	"time"

	svg "github.com/ajstarks/svgo"
)

type BuildParams struct {
	Data  []CalendarData
	Start time.Time
	End   time.Time
}

const (
	CellHeight         = 11
	CellWidth          = 11
	CellSpacing        = 2
	CellRoundingRadius = 1
	TopTextHeight      = CellHeight
)

func BuildSVG(params BuildParams) (string, error) {
	calendar := PrepareCalendar(params)
	buffer := bytes.Buffer{}
	image := svg.New(&buffer)

	cellsHeight := 7 * (CellHeight + CellSpacing)
	cellsWidth := (calendar.TotalWeeks) * (CellWidth + CellSpacing)
	height := cellsHeight + TopTextHeight + CellSpacing
	width := cellsWidth + CellWidth*2

	image.Start(width, height, `style="background:white"`)
	image.Style("text/css",
		`text { font-family: monospace; font-size: small; alignment-baseline: text-after-edge; }`,
		`rect { stroke: gray; stroke-width: 1; })`,
	)

	for _, point := range calendar.Points {
		image.Roundrect(
			point.Week*(CellWidth+CellSpacing),
			point.DayWeek*(CellHeight+CellSpacing)+TopTextHeight+CellSpacing,
			CellWidth,
			CellHeight,
			CellRoundingRadius,
			CellRoundingRadius,
			"fill:"+Color(point.Seconds, calendar.MaxSeconds),
		)
	}

	prev := ""
	for i, month := range calendar.Months {
		if month == prev {
			continue
		}

		prev = month

		image.Text(i*(CellWidth+CellSpacing), CellHeight+CellSpacing, month)
	}

	for i := time.Monday; i < time.Saturday; i += 2 {
		image.Text(cellsWidth, int(i+2)*(CellHeight+CellSpacing), i.String()[:3])
	}

	if params.Data == nil {
		image.Text(width, height, "Updating now",
			`style="text-anchor:end;font-size:large;fill:red;stroke:black;"`,
		)
	}

	image.End()

	return buffer.String(), nil
}
