package db

import "database/sql"

const EventTypeMessage int64 = 1
const EventTypeMessageEdited int64 = 2
const EventTypeReplyCallback int64 = 3

const MetricFlagFirstMessage = 1

type Metric struct {
	UserId    int64
	ChatId    int64
	EventType int64
	Content   string
	MediaType string
	MediaId   string
	AlbumId   string
	ReplyTo   int64
	Flags     int64
}

func WriteMetric(db *sql.DB, metric Metric) error {
	_, err := db.Exec(`INSERT INTO metrics (
		user_id,
		chat_id,
		event_type,
		content,
		media_type,
		media_id,
		album_id,
		reply_to,
		flags
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		metric.UserId,
		metric.ChatId,
		metric.EventType,
		metric.Content,
		metric.MediaType,
		metric.MediaId,
		metric.AlbumId,
		metric.ReplyTo,
		metric.Flags,
	)
	return err
}
