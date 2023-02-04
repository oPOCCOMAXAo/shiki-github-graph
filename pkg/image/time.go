package image

import "time"

func RoundDayStart(value time.Time, loc *time.Location) time.Time {
	value = value.In(loc)

	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, loc)
}

func RoundWeekStart(value time.Time, loc *time.Location) time.Time {
	value = value.In(loc)

	return time.Date(value.Year(), value.Month(), value.Day()-int(value.Weekday()), 0, 0, 0, 0, loc)
}

func RoundDayEnd(value time.Time, loc *time.Location) time.Time {
	value = value.In(loc)

	return time.Date(value.Year(), value.Month(), value.Day(), 23, 59, 59, 0, loc)
}

func GetDayPosition(value, start time.Time, loc *time.Location) PointPosition {
	value = value.In(loc)

	return PointPosition{
		Week:    int(value.Sub(start).Hours() / (24 * 7)),
		DayWeek: int(value.Weekday()),
	}
}
