package app

import (
	"log"
	"os"
)

func FileLogger(fname string) (*log.Logger, error) {
	file, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return nil, err
	}

	return log.New(file, log.Prefix(), log.Flags()), nil
}
