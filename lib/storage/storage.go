package storage

import (
	"database/sql"
	"log"

	"github.com/kdimonych/go_douuarss/lib/rss"
	_ "github.com/lib/pq"
)

type ChannelId int64
type ItemId int64

type Storage interface {
	Close() error
	InsertOrMergeChannel(channel *rss.Channel) (ChannelId, error)
	//TODO: AllChannels() ([]*rss.Channel, error)
}

type storage struct {
	DbURL string
	Db    *sql.DB
}

func NewStorage(dbURL string) (Storage, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &storage{
		DbURL: dbURL,
		Db:    db,
	}, nil
}

func (s *storage) Close() error {
	if s.Db != nil {
		if err := s.Db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (s *storage) insertOrReplaceChannel(channel *rss.Channel) (ChannelId, error) {
	query := `
		INSERT INTO channels (title, link, description, language, last_build_date)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (title) DO UPDATE SET
			link = EXCLUDED.link,
			description = EXCLUDED.description,
			language = EXCLUDED.language,
			last_build_date = EXCLUDED.last_build_date
		RETURNING id;
	`
	var id ChannelId
	err := s.Db.QueryRow(query, channel.Title, channel.Link, channel.Description, channel.Language, channel.LastBuildDate.Time).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) insertOrReplaceItem(item *rss.Item, channelId ChannelId) (ItemId, error) {
	query := `
		INSERT INTO items (title, link, description, pub_date, creator, channel_id, hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (hash) DO UPDATE SET
			link = EXCLUDED.link,
			description = EXCLUDED.description,
			creator = EXCLUDED.creator
		RETURNING id;
	`
	hash := item.Title + "_" + item.PubDate.String()
	var id ItemId
	err := s.Db.QueryRow(query, item.Title, item.Link, item.Description, item.PubDate.Time, item.Creator, channelId, hash).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) InsertOrMergeChannel(channel *rss.Channel) (ChannelId, error) {
	id, err := s.insertOrReplaceChannel(channel)
	if err != nil {
		return 0, err
	}

	for _, item := range channel.Items {
		_, err := s.insertOrReplaceItem(&item, id)
		if err != nil {
			log.Println("Unable to insert the item %w: %w", item, err)
			return 0, err
		}
	}
	return id, nil
}
