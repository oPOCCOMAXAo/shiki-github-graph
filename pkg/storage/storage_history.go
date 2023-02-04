package storage

import (
	"context"

	"github.com/pkg/errors"
)

func (s *Storage) SaveHistory(
	ctx context.Context,
	history []*History,
) error {
	if len(history) == 0 {
		return nil
	}

	err := s.db.WithContext(ctx).
		Save(&history).
		Error

	return errors.WithStack(err)
}

func (s *Storage) GetHistoryBroken(
	ctx context.Context,
	userID int64,
) ([]*History, error) {
	var res []*History

	err := s.db.WithContext(ctx).
		Where("episode_start = 0").
		Where("user_id = ?", userID).
		Find(&res).
		Error

	return res, errors.WithStack(err)
}

func (s *Storage) GetHistoryAnimes(
	ctx context.Context,
	userID int64,
	animesID []int32,
) ([]*History, error) {
	var res []*History

	err := s.db.WithContext(ctx).
		Where("anime_id IN (?)", animesID).
		Where("user_id = ?", userID).
		Find(&res).
		Error

	return res, errors.WithStack(err)
}
