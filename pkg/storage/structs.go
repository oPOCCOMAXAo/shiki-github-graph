package storage

type User struct {
	ID            int64  `gorm:"column:id;primaryKey"`
	Name          string `gorm:"column:name;index"`
	LastUpdatedAt int64  `gorm:"column:last_updated_at;not null;default:0;index"`
	RequestedAt   int64  `gorm:"column:requested_at;not null;default:0"`
	Full          *bool  `gorm:"column:full;not null;default:0"`
}

func (User) TableName() string {
	return "user"
}

type History struct {
	ID           int64 `gorm:"column:id;primaryKey"`
	CreatedAt    int64 `gorm:"column:created_at"`
	UserID       int64 `gorm:"column:user_id"`
	AnimeID      int32 `gorm:"column:anime_id"`
	EpisodeStart int16 `gorm:"column:episode_start"`
	EpisodeEnd   int16 `gorm:"column:episode_end"`
}

func (History) TableName() string {
	return "history"
}

type Anime struct {
	ID                     int64  `gorm:"column:id;primaryKey"`
	Name                   string `gorm:"column:name;not null"`
	Episodes               int    `gorm:"column:episodes;not null"`
	UpdatedAt              int64  `gorm:"column:updated_at;not null;default:0;index"`
	EpisodeDurationSeconds int    `gorm:"column:episode_duration_seconds;not null"`
}

func (Anime) TableName() string {
	return "anime"
}

type CalendarPoint struct {
	Time    int64 `gorm:"column:time"`
	Seconds int64 `gorm:"column:seconds"`
}
