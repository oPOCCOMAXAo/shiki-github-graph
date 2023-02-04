package shiki

import "time"

type Anime struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Episodes int    `json:"episodes"`
	Duration int    `json:"duration"`
}

type HistoryEntry struct {
	ID          int64             `json:"id"`
	CreatedAt   time.Time         `json:"created_at"`
	Description string            `json:"description"`
	Target      *Anime            `json:"target"`
	Watch       *WatchDescription `json:"-"`
}

type User struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}
