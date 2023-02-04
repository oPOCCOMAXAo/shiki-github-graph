package storage

import (
	"log"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/opoccomaxao/shiki-github-graph/pkg/app"
)

type Storage struct {
	db *gorm.DB
}

type Config struct {
	DB string `env:"DB,required"`
}

func New(config Config) (*Storage, error) {
	if config.DB == "" {
		return nil, errors.WithMessage(app.ErrParameterRequired, "DB")
	}

	res := Storage{}

	var err error

	res.db, err = gorm.Open(mysql.Open(config.DB), &gorm.Config{
		Logger: logger.New(log.Default(), logger.Config{
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &res, res.init()
}

func (s *Storage) init() error {
	return s.initMigrations()
}
