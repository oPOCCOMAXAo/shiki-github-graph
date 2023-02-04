package storage

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (s *Storage) GetUserByName(
	ctx context.Context,
	name string,
) (*User, error) {
	var res User

	err := s.db.WithContext(ctx).
		Where("name = ?", name).
		Take(&res).
		Error

	if err == nil {
		return &res, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, errors.WithStack(err)
}

func (s *Storage) SaveUser(
	ctx context.Context,
	user *User,
) error {
	omit := []string{}

	if user.LastUpdatedAt == 0 {
		omit = append(omit, "last_updated_at")
	}

	if user.RequestedAt == 0 {
		omit = append(omit, "requested_at")
	}

	if user.Name == "" {
		omit = append(omit, "name")
	}

	if user.Full == nil {
		omit = append(omit, "full")
	}

	err := s.db.WithContext(ctx).
		Omit(omit...).
		Save(user).
		Error

	return errors.WithStack(err)
}

func (s *Storage) GetUserToUpdate(ctx context.Context) (*User, error) {
	var res User

	now := time.Now().Unix()

	err := s.db.WithContext(ctx).
		Where("requested_at > last_updated_at").
		Where("requested_at > ?", now-LastValidRequested).
		Where("last_updated_at < ? OR full = 0", now-FirstValidUpdated).
		Order("full ASC").
		Order("last_updated_at ASC").
		Take(&res).
		Error

	if err == nil {
		return &res, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, errors.WithStack(err)
}
