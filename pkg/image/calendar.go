package image

import (
	"time"

	"github.com/opoccomaxao-go/generic-collection/gmath"
)

type CalendarData struct {
	Day     time.Time
	Seconds int64
}

type PointPosition struct {
	Week    int
	DayWeek int
}

type Point struct {
	PointPosition
	Seconds int64
	Label   string
}

type CalendarSVG struct {
	Points     []*Point
	Months     []string
	MaxSeconds int64
	TotalWeeks int
}

const FormatDate = "2006-01-02"

func PrepareCalendar(params BuildParams) CalendarSVG {
	mapDays := map[PointPosition]*Point{}
	location := time.UTC

	params.Start = RoundWeekStart(params.Start, location)
	params.End = RoundDayEnd(params.End, location)

	res := CalendarSVG{
		TotalWeeks: GetDayPosition(params.End, params.Start, location).Week + 1,
		MaxSeconds: 1,
	}
	res.Months = make([]string, res.TotalWeeks)

	for i := params.Start; !i.After(params.End); i = i.AddDate(0, 0, 1) {
		point := Point{
			PointPosition: GetDayPosition(i, params.Start, location),
			Seconds:       0,
			Label:         i.Format(FormatDate),
		}

		mapDays[point.PointPosition] = &point
		res.Points = append(res.Points, &point)

		if i.Day() < 7 {
			res.Months[point.Week] = i.AddDate(0, -1, 0).Month().String()[:3]
		} else {
			res.Months[point.Week] = i.Month().String()[:3]
		}
	}

	if res.Months[0] != res.Months[2] {
		res.Months[0] = res.Months[0][:1]
	}

	if res.Months[1] != res.Months[2] {
		res.Months[1] = res.Months[1][:1]
	}

	for _, data := range params.Data {
		pos := GetDayPosition(data.Day, params.Start, location)

		point := mapDays[pos]
		if point == nil {
			continue
		}

		point.Seconds += data.Seconds
	}

	res.MaxSeconds = 0
	for _, point := range res.Points {
		res.MaxSeconds = gmath.Max(point.Seconds, res.MaxSeconds)
	}

	return res
}
