package storage

import (
	"github.com/pkg/errors"
)

func (s *Storage) initMigrations() error {
	err := s.db.
		Migrator().
		AutoMigrate(
			&User{},
			&Anime{},
			&History{},
		)

	return errors.WithStack(err)
}
