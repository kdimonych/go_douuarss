package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

const InvalidFeedId FeedId = 0
const InvalidChannelId ChannelId = 0
const InvalidItemId ItemId = 0

type Storage interface {
	Close() error
	AddFeed(ctx context.Context, feed *Feed) (FeedId, error)
	GetFeedsOnly(ctx context.Context) ([]Feed, error)
	InsertOrMergeChannel(ctx context.Context, feedId FeedId, channel *Channel) (ChannelId, error)
	//TODO: GetAll() ([]Feed, error)
}

type storage struct {
	Db            *sql.DB
	channelHasher ChannelHasher
	itemHasher    ItemHasher
}

func (s *storage) Close() error {
	return s.Db.Close()
}

func (s *storage) AddFeed(ctx context.Context, feed *Feed) (FeedId, error) {
	if err := s.Db.PingContext(ctx); err != nil {
		return InvalidFeedId, err
	}

	id, err := s.addFeed(ctx, feed)
	if err != nil {
		return InvalidFeedId, err
	}

	for i := range feed.Channels {
		channel := &feed.Channels[i]
		_, err := s.InsertOrMergeChannel(ctx, id, channel)
		if err != nil {
			log.Println("Unable to insert the channel %w: %w", channel, err)
			return InvalidFeedId, err
		}
	}
	return id, nil
}

func (s *storage) GetFeedsOnly(ctx context.Context) ([]Feed, error) {
	if err := s.Db.PingContext(ctx); err != nil {
		return nil, err
	}

	query := `SELECT id, name, url FROM feeds;`
	rows, err := s.Db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []Feed
	for rows.Next() {
		var feed Feed
		if err := rows.Scan(&feed.Id, &feed.Name, &feed.Url); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (s *storage) InsertOrMergeChannel(ctx context.Context, feedId FeedId, channel *Channel) (ChannelId, error) {
	if err := s.Db.PingContext(ctx); err != nil {
		return InvalidChannelId, err
	}

	id, err := s.insertOrReplaceChannel(ctx, channel, feedId)
	if err != nil {
		return InvalidChannelId, err
	}

	for _, item := range channel.Items {
		_, err := s.insertOrReplaceItem(ctx, &item, id)
		if err != nil {
			log.Println("Unable to insert the item %w: %w", item, err)
			return InvalidChannelId, err
		}
	}
	return id, nil
}

// ==================== Storage Builder ====================
type StorageBuilder interface {
	WithChannelHasher(channelHasher ChannelHasher) StorageBuilder
	WithItemHasher(itemHasher ItemHasher) StorageBuilder
	WithDbConnectionFabric(dbConnectionFabric DbConnectionFabric) StorageBuilder
	Build(dbURL string) (Storage, error)
}

type storageBuilder struct {
	channelHasher      ChannelHasher
	itemHasher         ItemHasher
	dbConnectionFabric DbConnectionFabric
}

func NewStorageBuilder() StorageBuilder {
	return &storageBuilder{channelHasher: NewChannelHasher(),
		itemHasher:         NewItemHasher(),
		dbConnectionFabric: NewDbConnectionFabric(),
	}
}

func (b *storageBuilder) WithChannelHasher(channelHasher ChannelHasher) StorageBuilder {
	if channelHasher == nil {
		log.Println("Nil ChannelHasher provided. Using default ChannelHasher")
		channelHasher = NewChannelHasher()
	}
	b.channelHasher = channelHasher
	return b
}

func (b *storageBuilder) WithItemHasher(itemHasher ItemHasher) StorageBuilder {
	if itemHasher == nil {
		log.Println("Nil ItemHasher provided. Using default ItemHasher")
		itemHasher = NewItemHasher()
	}
	b.itemHasher = itemHasher
	return b
}

func (b *storageBuilder) WithDbConnectionFabric(dbConnectionFabric DbConnectionFabric) StorageBuilder {
	if dbConnectionFabric == nil {
		log.Println("Nil DbConnectionFabric provided. Using default DbConnectionFabric")
		dbConnectionFabric = NewDbConnectionFabric()
	}
	b.dbConnectionFabric = dbConnectionFabric
	return b
}

func (b *storageBuilder) Build(dbURL string) (Storage, error) {
	if dbURL == "" {
		return nil, errors.New("invalid database URL")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return &storage{
		Db:            db,
		channelHasher: b.channelHasher,
		itemHasher:    b.itemHasher,
	}, nil
}

// ==================== Private methods ====================

func (s *storage) addFeed(ctx context.Context, feed *Feed) (FeedId, error) {
	query := `
		INSERT INTO feeds (name, url)
		VALUES ($1, $2)
		ON CONFLICT (url) DO UPDATE SET
			name = EXCLUDED.name
		RETURNING id;
	`

	var id FeedId
	err := s.Db.QueryRowContext(ctx,
		query,
		feed.Name,
		feed.Url,
	).Scan(&id)
	if err != nil {
		return InvalidFeedId, &StorageError{
			Code:        ErrorUnableToInsertFeed,
			Description: "Feed url: " + feed.Url,
			Details:     err,
		}
	}

	return id, nil
}

func (s *storage) insertOrReplaceChannel(ctx context.Context, channel *Channel, feedId FeedId) (ChannelId, error) {
	query := `
		INSERT INTO channels (title, link, description, language, last_build_date, feed_id, hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (hash) DO UPDATE SET
			title = EXCLUDED.title,
			link = EXCLUDED.link,
			description = EXCLUDED.description,
			language = EXCLUDED.language,
			last_build_date = EXCLUDED.last_build_date
		RETURNING id;
	`

	channelHash := strconv.FormatInt(int64(feedId), 10) + "_" + s.channelHasher.Hash(channel)

	var id ChannelId
	err := s.Db.QueryRowContext(ctx,
		query,
		channel.Title,
		channel.Link,
		channel.Description,
		channel.Language,
		channel.LastBuildDate,
		feedId,
		channelHash,
	).Scan(&id)
	if err != nil {
		return InvalidChannelId, &StorageError{
			Code:        ErrorUnableToInsertChannel,
			Description: "channelHash: " + channelHash,
			Details:     err}
	}

	return id, nil
}

func (s *storage) insertOrReplaceItem(ctx context.Context, item *Item, channelId ChannelId) (ItemId, error) {
	query := `
		INSERT INTO items (title, link, description, pub_date, creator, channel_id, hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (hash) DO UPDATE SET
			title = EXCLUDED.title,
			link = EXCLUDED.link,
			description = EXCLUDED.description,
			pub_date = EXCLUDED.pub_date,
			creator = EXCLUDED.creator
		RETURNING id;
	`

	itemHash := strconv.FormatInt(int64(channelId), 10) + "_" + s.itemHasher.Hash(item)

	var id ItemId
	err := s.Db.QueryRowContext(ctx,
		query,
		item.Title,
		item.Link,
		item.Description,
		item.PubDate,
		item.Creator,
		channelId,
		itemHash).Scan(&id)
	if err != nil {
		return InvalidItemId, &StorageError{
			Code:        ErrorUnableToInsertItem,
			Description: "itemHash: " + itemHash,
			Details:     err,
		}
	}

	return id, nil
}
