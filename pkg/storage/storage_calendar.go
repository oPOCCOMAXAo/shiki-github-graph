package storage

import (
	"context"

	"github.com/pkg/errors"
)

func (s *Storage) GetUserCalendar(
	ctx context.Context,
	userID int64,
) ([]*CalendarPoint, error) {
	var res []*CalendarPoint

	err := s.db.WithContext(ctx).
		Select(
			"h.created_at AS time",
			"(h.episode_end - h.episode_start + 1) * a.episode_duration_seconds AS seconds",
		).
		Table("history h").
		Joins("JOIN anime a ON a.id = h.anime_id").
		Where("h.user_id = ?", userID).
		Find(&res).
		Error

	return res, errors.WithStack(err)
}

func (s *Storage) IsCalendarUnfinished(
	ctx context.Context,
	userID int64,
) (bool, error) {
	var total int64

	err := s.db.WithContext(ctx).
		Table("history h").
		Joins("LEFT JOIN anime a ON a.id = h.anime_id").
		Where("a.id IS NULL").
		Count(&total).
		Error

	return total == 0, errors.WithStack(err)
}
