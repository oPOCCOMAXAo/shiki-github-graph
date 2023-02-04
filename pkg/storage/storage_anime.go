package storage

import (
	"context"

	"github.com/pkg/errors"
)

func (s *Storage) SaveAnime(ctx context.Context, anime []*Anime) error {
	if len(anime) == 0 {
		return nil
	}

	err := s.db.WithContext(ctx).
		Save(&anime).
		Error

	return errors.WithStack(err)
}

func (s *Storage) GetAnimesUnstored(
	ctx context.Context,
	limit int,
) ([]int32, error) {
	var res []int32

	if limit <= 0 {
		limit = 50
	}

	err := s.db.WithContext(ctx).
		Table("history h").
		Joins("LEFT JOIN anime a ON a.id = h.anime_id").
		Where("a.id IS NULL").
		Group("h.anime_id").
		Order("h.id DESC").
		Limit(limit).
		Pluck("h.anime_id", &res).
		Error

	return res, errors.WithStack(err)
}
